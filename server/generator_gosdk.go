package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	_ "github.com/golang/protobuf/protoc-gen-go/grpc"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jhump/protoreflect/desc"
)

func lookupFieldByName(kind *SchemaKind, name string) *SchemaField {
	for _, field := range kind.Fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}

func generateGoCode(fileDesc *desc.FileDescriptor, schema *Schema) (string, error) {
	g := generator.New()

	packageName := schema.Name
	protoFileName := fmt.Sprintf("%s.proto", schema.Name)
	parameter := "plugins=grpc"

	timestampFileDescriptor, err := desc.LoadFileDescriptor("google/protobuf/timestamp.proto")
	if err != nil {
		return "", err
	}
	timestampFileProto := timestampFileDescriptor.AsFileDescriptorProto()

	str := packageName
	strName := protoFileName
	fileDescProto := fileDesc.AsFileDescriptorProto()
	fileDescProto.Options = &dpb.FileOptions{
		GoPackage: &str,
	}
	fileDescProto.Name = &strName

	g.Request.Parameter = &parameter
	g.Request.FileToGenerate = append(
		g.Request.FileToGenerate,
		protoFileName,
	)
	g.Request.ProtoFile = append(
		g.Request.ProtoFile,
		fileDescProto,
	)
	g.Request.ProtoFile = append(
		g.Request.ProtoFile,
		timestampFileProto,
	)

	genFiles := make(map[string]*generator.FileDescriptor)
	genFiles[protoFileName] = &generator.FileDescriptor{
		FileDescriptorProto: fileDescProto,
	}
	genFiles["google/protobuf/timestamp.proto"] = &generator.FileDescriptor{
		FileDescriptorProto: timestampFileProto,
	}

	g.CommandLineParameters(g.Request.GetParameter())
	g.SetFiles(genFiles, protoFileName)
	g.WrapTypes()
	g.SetPackageNames()
	g.BuildTypeNameMap()
	g.GenerateAllFiles()

	standardCode := *g.Response.File[0].Content

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	tmplCode, err := ioutil.ReadFile(path.Join(dir, "generator_gosdk_template.gotxt"))
	if err != nil {
		return "", err
	}

	tmpl := template.New("gosdk")
	tmpl.Funcs(map[string]interface{}{
		"getschema": func() *Schema {
			return schema
		},
		"getschemakind": func(kindName string) *SchemaKind {
			return schema.Kinds[kindName]
		},
		"camelcase": func(str string) string {
			return generator.CamelCase(str)
		},
		"getfieldforindex": func(kindName string, targetIndex *SchemaIndex) *SchemaField {
			return lookupFieldByName(schema.Kinds[kindName], targetIndex.GetField())
		},
		"isinmemoryindex": func(index *SchemaIndex) bool {
			return index.Type == SchemaIndexType_memory
		},
		"isfieldindex": func(index *SchemaIndex) bool {
			switch index.Value.(type) {
			case *SchemaIndex_Field:
				return true
			default:
				return false
			}
		},
		"iscomputedindex": func(index *SchemaIndex) bool {
			switch index.Value.(type) {
			case *SchemaIndex_Computed:
				return true
			default:
				return false
			}
		},
		"iscomputedfnv64aindex": func(index *SchemaIndex) bool {
			switch index.Value.(type) {
			case *SchemaIndex_Computed:
				computed := index.GetComputed()
				switch computed.Algorithm.(type) {
				case *SchemaComputedIndex_Fnv64A:
					return true
				default:
					return false
				}
			default:
				return false
			}
		},
		"iscomputedfnv64apairindex": func(index *SchemaIndex) bool {
			switch index.Value.(type) {
			case *SchemaIndex_Computed:
				computed := index.GetComputed()
				switch computed.Algorithm.(type) {
				case *SchemaComputedIndex_Fnv64APair:
					return true
				default:
					return false
				}
			default:
				return false
			}
		},
		"getcomputedfnv64afieldforindex": func(kindName string, index *SchemaIndex) *SchemaField {
			computed := index.GetComputed()
			return lookupFieldByName(schema.Kinds[kindName], computed.GetFnv64A().Field)
		},
		"getcomputedfnv64apairfield1forindex": func(kindName string, index *SchemaIndex) *SchemaField {
			computed := index.GetComputed()
			return lookupFieldByName(schema.Kinds[kindName], computed.GetFnv64APair().Field1)
		},
		"getcomputedfnv64apairfield2forindex": func(kindName string, index *SchemaIndex) *SchemaField {
			computed := index.GetComputed()
			return lookupFieldByName(schema.Kinds[kindName], computed.GetFnv64APair().Field2)
		},
	})
	_, err = tmpl.Parse(string(tmplCode))
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, schema)
	if err != nil {
		return "", err
	}

	// replace
	imports := []string{
		"\"io\"",
		"\"sync\"",
		"\"strings\"",
		"\"time\"",
		"\"encoding/binary\"",
		"\"google.golang.org/grpc/status\"",
		"\"google.golang.org/grpc/codes\"",
		"\"hash/fnv\"",
	}
	if !strings.Contains(standardCode, "github.com/golang/protobuf/ptypes/timestamp") {
		imports = append(
			imports,
			"timestamp \"github.com/golang/protobuf/ptypes/timestamp\"",
		)
	}
	standardCode = strings.Replace(standardCode, "import (", fmt.Sprintf("import (%s", strings.Join(imports, "\n    ")), 1)

	standardCode = fmt.Sprintf("%s\n%s", standardCode, tpl.String())
	return standardCode, nil
	/*

			configstoreImplementationCode := `
		type Configstore struct {
			conn *grpc.ClientConn
			transactions []*MetaTransactionRecord

		`
			for kindName := range schema.Kinds {
				configstoreImplementationCode = strings.Join([]string{configstoreImplementationCode, `
			pending`, kindName, ` []`, kindName, `
			`, kindName, `s *`, kindName, `Store
		`}, "")
			}
			configstoreImplementationCode = strings.Join([]string{configstoreImplementationCode, `
			sync.RWMutex storeMutex
			sync.RWMutex pendingMutex
			sync.RWMutex transactionsMutex
		}

		func (configstore *Configstore) startListeningForTransactions() error {
			client := NewConfigstoreMetaServiceClient(configstore.Conn)
			watcher, err := client.Watch(ctx, &WatchTransactionsRequest{})
			if err != nil {
				return err
			}
			go func() {
				for {
					resp, err := watcher.Recv()
					if err == io.EOF {
						break
					}
					if err == nil || status.Code(err) == codes.OK {
						if resp == nil {
							time.Sleep(1 * time.Second)
							continue
						}
						if resp.Type == WatchEventType_Created {
							configstore.transactionsMutex.Lock()
							configstore.transactions = append(configstore.transactions, resp)
							sort.Slice(configstore.transactions[:], func(i, j int) bool {
								if (configstore.transactions[i].Seconds < configstore.transactions[j].Seconds) {
									return true
								}
								if (configstore.transactions[i].Nanos < configstore.transactions[j].Nanos) {
									return true
								}
								return false
							})
							configstore.transactionsMutex.Unlock()

							configstore.rescanPendingTransactions()
						}
					} else if status.Code(err) == codes.Unavailable {
						// Retry the Watch request itself
						watcher, err = client.Watch(ctx, &WatchTransactionsRequest{})
						if err != nil {
							time.Sleep(30 * time.Second)
							continue
						} else {
							// Connection re-established, loop back around again to receive updates.
						}
					} else {
						time.Sleep(1 * time.Second)
						continue
					}
				}
			}()
			return nil
		}

		func (configstore *Configstore) rescanPendingTransactions() error {
			configstore.transactionsMutex.RLock()
			defer configstore.transactionsMutex.RUnlock()

			// transactions are in order from earliest to latest

			for _, transaction := range configstore.transactions {
				if len(transaction.MutatedKeys) == 0 {

				}

				// check to see if this transaction has all the mutated entities it needs
				hasAllMutatedKeys := true

			}
		}

		func NewConfigstore(ctx context.Context, conn *grpc.ClientConn) *Configstore {
			configstore := &Configstore{
				Conn: conn,
			}
		`}, "")
			for kindName := range schema.Kinds {
				configstoreImplementationCode = strings.Join([]string{configstoreImplementationCode, `
			configstore.`, kindName, `s = New`, kindName, `Store(ctx, New`, kindName, `Client(conn), configstore)`}, "")
			}
			configstoreImplementationCode = strings.Join([]string{configstoreImplementationCode, `
			configstore.startListeningForTransactions()
		}
		`}, "")

			// extendedCode = fmt.Sprintf("%s\n%s", extendedCode, configstoreImplementationCode)

			for kindName := range schema.Kinds {
				indexStores := ""
				indexFuncDecls := ""
				indexFuncDefs := ""
				indexStoresInit := ""
				indexStoresRemove := ""
				indexStoresUpdate := ""

				for _, index := range schema.Kinds[kindName].Indexes {
					if index.Type == SchemaIndexType_memory {
						indexKeyType := ""
						indexOriginalType := ""
						indexSerializeToType := "key"

						switch index.Value.(type) {
						case *SchemaIndex_Field:
							field := lookupFieldByName(schema.Kinds[kindName], index.GetField())
							if field == nil {
								continue
							}

							switch field.Type {
							case ValueType_double:
								indexKeyType = "float64"
								indexOriginalType = "float64"
								break
							case ValueType_int64:
								indexKeyType = "int64"
								indexOriginalType = "int64"
								break
							case ValueType_uint64:
								indexKeyType = "uint64"
								indexOriginalType = "uint64"
								break
							case ValueType_timestamp:
								indexKeyType = "string"
								indexOriginalType = "*timestamp.Timestamp"
								indexSerializeToType = "SerializeTimestamp(key)"
								break
							case ValueType_bytes:
								indexKeyType = "string"
								indexOriginalType = "[]byte"
								indexSerializeToType = "string(key)"
								break
							case ValueType_key:
								indexKeyType = "string"
								indexOriginalType = "*Key"
								indexSerializeToType = "SerializeKey(key)"
								break
							case ValueType_boolean:
								indexKeyType = "bool"
								indexOriginalType = "bool"
								break
							case ValueType_string:
								indexKeyType = "string"
								indexOriginalType = "string"
								break
							default:
								return "", fmt.Errorf("unexpected field type for index")
							}

							indexStoresUpdate = fmt.Sprintf(
								`%s
				{
					key := newEntity.%s
					idx := %s
					ref.indexstore_%s[idx] = newEntity
				}
			`,
								indexStoresUpdate,
								generator.CamelCase(field.Name),
								indexSerializeToType,
								index.Name,
							)
							indexStoresRemove = fmt.Sprintf(
								`%s
				{
					key := oldEntity.%s
					idx := %s
					delete(ref.indexstore_%s, idx)
				}
			`,
								indexStoresRemove,
								generator.CamelCase(field.Name),
								indexSerializeToType,
								index.Name,
							)
							indexFuncDefs = fmt.Sprintf(
								`%s
			func (c *<ENTITY>ImplStore) GetBy%s(key %s) *<ENTITY> {
				c.RLock()
				defer c.RUnlock()
				return c.indexstore_%s[%s]
			}

			func (c *<ENTITY>ImplStore) GetAndCheckBy%s(key %s) (*<ENTITY>, bool) {
				c.RLock()
				defer c.RUnlock()
				val, ok := c.indexstore_%s[%s]
				return val, ok
			}`,
								indexFuncDefs,
								index.Name,
								indexOriginalType,
								index.Name,
								indexSerializeToType,
								index.Name,
								indexOriginalType,
								index.Name,
								indexSerializeToType,
							)
							break
						case *SchemaIndex_Computed:
							computed := index.GetComputed()
							switch computed.Algorithm.(type) {
							case *SchemaComputedIndex_Fnv64A:
								field := lookupFieldByName(schema.Kinds[kindName], computed.GetFnv64A().Field)
								if field == nil {
									continue
								}

								switch field.Type {
								case ValueType_double:
								case ValueType_int64:
								case ValueType_uint64:
								case ValueType_timestamp:
								case ValueType_bytes:
								case ValueType_key:
								case ValueType_boolean:
									return "", fmt.Errorf("unsupported field type for fnv64a index")
								case ValueType_string:
									indexKeyType = "uint64"
									indexOriginalType = "uint64"
									indexSerializeToType = "Fnv64a(key)"
									break
								default:
									return "", fmt.Errorf("unexpected field type for index")
								}

								indexStoresUpdate = fmt.Sprintf(
									`%s
					{
						key := newEntity.%s
						idx := %s
						ref.indexstore_%s[idx] = newEntity
					}
				`,
									indexStoresUpdate,
									generator.CamelCase(field.Name),
									indexSerializeToType,
									index.Name,
								)
								indexStoresRemove = fmt.Sprintf(
									`%s
					{
						key := oldEntity.%s
						idx := %s
						delete(ref.indexstore_%s, idx)
					}
				`,
									indexStoresRemove,
									generator.CamelCase(field.Name),
									indexSerializeToType,
									index.Name,
								)
								indexFuncDefs = fmt.Sprintf(
									`%s
				func (c *<ENTITY>ImplStore) GetBy%s(key %s) *<ENTITY> {
					c.RLock()
					defer c.RUnlock()
					return c.indexstore_%s[key]
				}

				func (c *<ENTITY>ImplStore) GetAndCheckBy%s(key %s) (*<ENTITY>, bool) {
					c.RLock()
					defer c.RUnlock()
					val, ok := c.indexstore_%s[key]
					return val, ok
				}`,
									indexFuncDefs,
									index.Name,
									indexOriginalType,
									index.Name,
									index.Name,
									indexOriginalType,
									index.Name,
								)
								break
							case *SchemaComputedIndex_Fnv64APair:
								field1 := lookupFieldByName(schema.Kinds[kindName], computed.GetFnv64APair().Field1)
								if field1 == nil {
									continue
								}
								field2 := lookupFieldByName(schema.Kinds[kindName], computed.GetFnv64APair().Field2)
								if field2 == nil {
									continue
								}

								if field1.Type == ValueType_uint64 && field2.Type == ValueType_uint64 {
									indexKeyType = "uint64"
									indexOriginalType = "uint64"
									indexSerializeToType = "Fnv64aPair(field1, field2)"
								} else {
									return "", fmt.Errorf("unsupported field type for fnv64a pair index")
								}

								indexStoresUpdate = fmt.Sprintf(
									`%s
					{
						field1 := newEntity.%s
						field2 := newEntity.%s
						idx := %s
						ref.indexstore_%s[idx] = newEntity
					}
				`,
									indexStoresUpdate,
									generator.CamelCase(field1.Name),
									generator.CamelCase(field2.Name),
									indexSerializeToType,
									index.Name,
								)
								indexStoresRemove = fmt.Sprintf(
									`%s
					{
						field1 := oldEntity.%s
						field2 := oldEntity.%s
						idx := %s
						delete(ref.indexstore_%s, idx)
					}
				`,
									indexStoresRemove,
									generator.CamelCase(field1.Name),
									generator.CamelCase(field2.Name),
									indexSerializeToType,
									index.Name,
								)
								indexFuncDefs = fmt.Sprintf(
									`%s
				func (c *<ENTITY>ImplStore) GetBy%s(key %s) *<ENTITY> {
					c.RLock()
					defer c.RUnlock()
					return c.indexstore_%s[key]
				}

				func (c *<ENTITY>ImplStore) GetAndCheckBy%s(key %s) (*<ENTITY>, bool) {
					c.RLock()
					defer c.RUnlock()
					val, ok := c.indexstore_%s[key]
					return val, ok
				}`,
									indexFuncDefs,
									index.Name,
									indexOriginalType,
									index.Name,
									index.Name,
									indexOriginalType,
									index.Name,
								)
								break
							default:
								return "", fmt.Errorf("unexpected computed type for index")
							}
						default:
							return "", fmt.Errorf("unexpected index type for index")
						}

						indexStores = fmt.Sprintf(
							"%s\n  indexstore_%s map[%s]*<ENTITY>",
							indexStores,
							index.Name,
							indexKeyType,
						)
						indexFuncDecls = fmt.Sprintf(
							"%s\n  GetBy%s(key %s) *<ENTITY>\n  GetAndCheckBy%s(key %s) (*<ENTITY>, bool)",
							indexFuncDecls,
							index.Name,
							indexOriginalType,
							index.Name,
							indexOriginalType,
						)
						indexStoresInit = fmt.Sprintf(
							"%s\n  indexstore_%s: make(map[%s]*<ENTITY>),",
							indexStoresInit,
							index.Name,
							indexKeyType,
						)
					}
				}

				if indexStoresUpdate != "" {
					indexStoresUpdate = fmt.Sprintf("newEntity := resp.Entity\n%s", indexStoresUpdate)
				}
				if indexStoresRemove != "" {
					indexStoresRemove = fmt.Sprintf(`
			oldEntity, ok := ref.store[s]
			if ok {
				%s
			}`, indexStoresRemove)
				}

				m := map[string]string{
					"<INDEXSTORES>":       indexStores,
					"<INDEXFUNCDECLS>":    indexFuncDecls,
					"<INDEXFUNCDEFS>":     indexFuncDefs,
					"<INDEXSTORESINIT>":   indexStoresInit,
					"<INDEXSTORESUPDATE>": indexStoresUpdate,
					"<INDEXSTORESREMOVE>": indexStoresRemove,
				}

				v := extendedTemplateCode
				for k, vv := range m {
					v = strings.Replace(v, k, vv, -1)
				}

				extendedCode = fmt.Sprintf(
					"%s\n%s",
					extendedCode,
					strings.Replace(v, "<ENTITY>", kindName, -1),
				)
			}

			return extendedCode, nil
	*/
}
