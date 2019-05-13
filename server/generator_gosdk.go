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
			if computed.GetFnv64A().Field == "key" {
				return &SchemaField{
					Id: 1,
					Name: "Key",
					Type: ValueType_key,
				}
			}
			return lookupFieldByName(schema.Kinds[kindName], computed.GetFnv64A().Field)
		},
		"getcomputedfnv64apairfield1forindex": func(kindName string, index *SchemaIndex) *SchemaField {
			computed := index.GetComputed()
			if computed.GetFnv64APair().Field1 == "key" {
				return &SchemaField{
					Id: 1,
					Name: "Key",
					Type: ValueType_key,
				}
			}
			return lookupFieldByName(schema.Kinds[kindName], computed.GetFnv64APair().Field1)
		},
		"getcomputedfnv64apairfield2forindex": func(kindName string, index *SchemaIndex) *SchemaField {
			computed := index.GetComputed()
			if computed.GetFnv64APair().Field2 == "key" {
				return &SchemaField{
					Id: 1,
					Name: "Key",
					Type: ValueType_key,
				}
			}
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
		"\"os\"",
		"\"log\"",
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
}
