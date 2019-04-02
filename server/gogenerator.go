package main

import (
	"fmt"
	"strings"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	_ "github.com/golang/protobuf/protoc-gen-go/grpc"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jhump/protoreflect/desc"
)

const (
	extendedOnceCode = `
func Fnv64a(val string) uint64 {
	hash := fnv.New64a()
	hash.Write(([]byte)(val))
	return hash.Sum64()
}

func Fnv64aPair(a uint64, b uint64) uint64 {
	hash := fnv.New64a()
	if a > b {
		tmp := a
		a = b
		b = tmp
	}
	key := make([]byte, 16)
	binary.LittleEndian.PutUint64(key, a)
	binary.LittleEndian.PutUint64(key[8:], b)
	hash.Write(key)
	return hash.Sum64()
}

func CreateTopLevelKey(partitionId *PartitionId, pathElement *PathElement) *Key {
	return &Key{
		PartitionId: partitionId,
		Path: []*PathElement{pathElement},
	}
}

func CreateIncompleteTopLevelKey(partitionId *PartitionId, kind string) *Key {
	return &Key{
		PartitionId: partitionId,
		Path: []*PathElement{
			&PathElement{
				Kind: kind,
			},
		},
	}
}

func CreateDescendantKey(parent *Key, pathElement *PathElement) *Key {
	newKey := &Key{
		PartitionId: &PartitionId{
			Namespace: parent.PartitionId.Namespace,
		},
		Path: nil,
	}
	for _, elem := range parent.Path {
		switch elem.IdType.(type) {
		case *PathElement_Id:
			newKey.Path = append(newKey.Path, &PathElement{
				Kind: elem.Kind,
				IdType: &PathElement_Id{
					Id: elem.GetId(),
				},
			})
			break
		case *PathElement_Name:
			newKey.Path = append(newKey.Path, &PathElement{
				Kind: elem.Kind,
				IdType: &PathElement_Name{
					Name: elem.GetName(),
				},
			})
			break
		}
	}
	switch pathElement.IdType.(type) {
	case *PathElement_Id:
		newKey.Path = append(newKey.Path, &PathElement{
			Kind: pathElement.Kind,
			IdType: &PathElement_Id{
				Id: pathElement.GetId(),
			},
		})
		break
	case *PathElement_Name:
		newKey.Path = append(newKey.Path, &PathElement{
			Kind: pathElement.Kind,
			IdType: &PathElement_Name{
				Name: pathElement.GetName(),
			},
		})
		break
	}
	return newKey
}

func SerializeTimestamp(ts *timestamp.Timestamp) string {
	if ts == nil {
		return ""
	}

	return ts.String()
}

func SerializeKey(key *Key) string {
	if key == nil {
		return ""
	}

	var elements []string
	for _, pathElement := range key.Path {
		if _, ok := pathElement.IdType.(*PathElement_Id); ok {
			elements = append(elements, fmt.Sprintf("id=%d", pathElement.GetId()))
		} else {
			elements = append(elements, fmt.Sprintf("name=%s", pathElement.GetName()))
		}
	}
	return fmt.Sprintf("ns=%s|%s", key.PartitionId.Namespace, strings.Join(elements, "|"))
}

func CompareKeys(a *Key, b *Key) bool {
	return SerializeKey(a) == SerializeKey(b)
}
`
	extendedTemplateCode = `
func CreateTopLevel_<ENTITY>_NameKey(partitionId *PartitionId, name string) *Key {
	return &Key{
		PartitionId: partitionId,
		Path: []*PathElement{
			&PathElement{
				Kind: "<ENTITY>",
				IdType: &PathElement_Name{
					Name: name,
				},
			},
		},
	}
}

func CreateTopLevel_<ENTITY>_IdKey(partitionId *PartitionId, id int64) *Key {
	return &Key{
		PartitionId: partitionId,
		Path: []*PathElement{
			&PathElement{
				Kind: "<ENTITY>",
				IdType: &PathElement_Id{
				  Id: id,
				},
			},
		},
	}
}

func CreateTopLevel_<ENTITY>_IncompleteKey(partitionId *PartitionId) *Key {
	return &Key{
		PartitionId: partitionId,
		Path: []*PathElement{
			&PathElement{
				Kind: "<ENTITY>",
				IdType: nil,
			},
		},
	}
}
	
type <ENTITY>ImplStore struct {
	sync.RWMutex
	client <ENTITY>ServiceClient
	store map[string]*<ENTITY>
	<INDEXSTORES>
}

type <ENTITY>Store interface {
	Create(ctx context.Context, entity *<ENTITY>) (*<ENTITY>, error)
	Update(ctx context.Context, entity *<ENTITY>) (*<ENTITY>, error)
	Delete(ctx context.Context, key *Key) (*<ENTITY>, error)
	GetAndCheck(key *Key) (*<ENTITY>, bool)
	Get(key *Key) *<ENTITY>
	GetKeys() []*Key
	<INDEXFUNCDECLS>
}

<INDEXFUNCDEFS>

func New<ENTITY>Store(ctx context.Context, client <ENTITY>ServiceClient) (<ENTITY>Store, error) {
	ref := &<ENTITY>ImplStore{
		client:   client,
		store:    make(map[string]*<ENTITY>),
		<INDEXSTORESINIT>
	}
	fmt.Printf("connecting to <ENTITY> store...\n")
	watcher, err := ref.client.Watch(ctx, &Watch<ENTITY>Request{})
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			resp, err := watcher.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				// TODO: Pass this to somewhere
				fmt.Printf("%v", err)
			}
			if resp.Type == WatchEventType_Created ||
			resp.Type == WatchEventType_Updated {
				ref.Lock()
				s := SerializeKey(resp.Entity.Key)
				<INDEXSTORESREMOVE>
				<INDEXSTORESUPDATE>
				ref.store[s] = resp.Entity
				ref.Unlock()
			} else if resp.Type == WatchEventType_Deleted {
				ref.Lock()
				s := SerializeKey(resp.Entity.Key)
				<INDEXSTORESREMOVE>
				delete(ref.store, s)
				ref.Unlock()
			}
		}
	}()
	fmt.Printf("listing <ENTITY> store...\n")
	req := &List<ENTITY>Request{
		Start: nil,
		Limit: 0,
	}
	for {
		resp, err := ref.client.List(ctx, req)
		if err != nil {
			return nil, err
		}
		for _, entity := range resp.Entities {
			ref.Lock()
			fmt.Printf("storing entity with key '%s'...\n", SerializeKey(entity.Key))
			ref.store[SerializeKey(entity.Key)] = entity
			ref.Unlock()
		}
		if !resp.MoreResults {
			break
		}
		req.Start = resp.Next
	}
	fmt.Printf("returning new <ENTITY> store...\n")
	return ref, nil
}

func (c *<ENTITY>ImplStore) GetKeys() []*Key {
	keys := make([]*Key, len(c.store))
	i := 0
	for _, entity := range c.store {
		keys[i] = entity.Key
		i += 1
	}
	return keys
}

func (c *<ENTITY>ImplStore) Get(key *Key) *<ENTITY> {
	c.RLock()
	defer c.RUnlock()
	return c.store[SerializeKey(key)]
}

func (c *<ENTITY>ImplStore) GetAndCheck(key *Key) (*<ENTITY>, bool) {
	c.RLock()
	defer c.RUnlock()
	fmt.Printf("trying to retrieve entity with key '%s'...\n", SerializeKey(key))
	val, ok := c.store[SerializeKey(key)]
	return val, ok
}

func (ref *<ENTITY>ImplStore) Create(ctx context.Context, entity *<ENTITY>) (*<ENTITY>, error) {
	resp, err := ref.client.Create(ctx, &Create<ENTITY>Request{
		Entity: entity,
	})
	if err != nil {
		return nil, err
	}
	s := SerializeKey(resp.Entity.Key)
	ref.Lock()
	<INDEXSTORESUPDATE>
	ref.store[s] = resp.Entity
	ref.Unlock()
	return resp.Entity, nil
}

func (ref *<ENTITY>ImplStore) Update(ctx context.Context, entity *<ENTITY>) (*<ENTITY>, error) {
	resp, err := ref.client.Update(ctx, &Update<ENTITY>Request{
		Entity: entity,
	})
	if err != nil {
		return nil, err
	}
	s := SerializeKey(resp.Entity.Key)
	ref.Lock()
	<INDEXSTORESREMOVE>
	<INDEXSTORESUPDATE>
	ref.store[s] = resp.Entity
	ref.Unlock()
	return resp.Entity, nil
}

func (ref *<ENTITY>ImplStore) Delete(ctx context.Context, key *Key) (*<ENTITY>, error) {
	resp, err := ref.client.Delete(ctx, &Delete<ENTITY>Request{
		Key: key,
	})
	if err != nil {
		return nil, err
	}
	s := SerializeKey(resp.Entity.Key)
	ref.Lock()
	<INDEXSTORESREMOVE>
	delete(ref.store, s)
	ref.Unlock()
	return resp.Entity, nil
}

`
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

	// Now add our automatically synchronising store code
	extendedCode := standardCode
	extendedCode = fmt.Sprintf("%s\n%s", standardCode, extendedOnceCode)
	if strings.Contains(extendedCode, "github.com/golang/protobuf/ptypes/timestamp") {
		extendedCode = strings.Replace(extendedCode, "import (", "import (\n    \"io\"\n    \"sync\"\n    \"strings\"\n    \"encoding/binary\"\n    \"hash/fnv\"", 1)
	} else {
		extendedCode = strings.Replace(extendedCode, "import (", "import (\n    \"io\"\n    \"sync\"\n    \"strings\"\n    \"encoding/binary\"\n    \"hash/fnv\"\n    timestamp \"github.com/golang/protobuf/ptypes/timestamp\"", 1)
	}
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
}
