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

func SerializeKey(key *Key) string {
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
}

type <ENTITY>Store interface {
	Create(ctx context.Context, entity *<ENTITY>) (*<ENTITY>, error)
	Update(ctx context.Context, entity *<ENTITY>) (*<ENTITY>, error)
	Delete(ctx context.Context, key *Key) (*<ENTITY>, error)
	GetAndCheck(key *Key) (*<ENTITY>, bool)
	Get(key *Key) *<ENTITY>
	GetKeys() []*Key
}

func New<ENTITY>Store(ctx context.Context, client <ENTITY>ServiceClient) (<ENTITY>Store, error) {
	ref := &<ENTITY>ImplStore{
		client: client,
		store:  make(map[string]*<ENTITY>),
	}
	watcher, err := ref.client.Watch(ctx, &Watch<ENTITY>Request{})
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			change, err := watcher.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				// TODO: Pass this to somewhere
				fmt.Printf("%v", err)
			}
			if change.Type == WatchEventType_Created ||
				change.Type == WatchEventType_Updated {
				ref.Lock()
				ref.store[SerializeKey(change.Entity.Key)] = change.Entity
				ref.Unlock()
			} else if change.Type == WatchEventType_Deleted {
				ref.Lock()
				delete(ref.store, SerializeKey(change.Entity.Key))
				ref.Unlock()
			}
		}
	}()
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
			ref.store[SerializeKey(entity.Key)] = entity
			ref.Unlock()
		}
		if !resp.MoreResults {
			break
		}
		req.Start = resp.Next
	}
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
	val, ok := c.store[SerializeKey(key)]
	return val, ok
}

func (c *<ENTITY>ImplStore) Create(ctx context.Context, entity *<ENTITY>) (*<ENTITY>, error) {
	resp, err := c.client.Create(ctx, &Create<ENTITY>Request{
		Entity: entity,
	})
	if err != nil {
		return nil, err
	}
	c.Lock()
	c.store[SerializeKey(resp.Entity.Key)] = resp.Entity
	c.Unlock()
	return resp.Entity, nil
}

func (c *<ENTITY>ImplStore) Update(ctx context.Context, entity *<ENTITY>) (*<ENTITY>, error) {
	resp, err := c.client.Update(ctx, &Update<ENTITY>Request{
		Entity: entity,
	})
	if err != nil {
		return nil, err
	}
	c.Lock()
	c.store[SerializeKey(resp.Entity.Key)] = resp.Entity
	c.Unlock()
	return resp.Entity, nil
}

func (c *<ENTITY>ImplStore) Delete(ctx context.Context, key *Key) (*<ENTITY>, error) {
	resp, err := c.client.Delete(ctx, &Delete<ENTITY>Request{
		Key: key,
	})
	if err != nil {
		return nil, err
	}
	c.Lock()
	delete(c.store, SerializeKey(resp.Entity.Key))
	c.Unlock()
	return resp.Entity, nil
}

`
)

func generateGoCode(fileDesc *desc.FileDescriptor, schema *configstoreSchema) (string, error) {
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
	extendedCode = strings.Replace(extendedCode, "import (", "import (\n    \"io\"\n    \"sync\"\n    \"strings\"", 1)
	for kindName := range schema.Kinds {
		extendedCode = fmt.Sprintf(
			"%s\n%s",
			extendedCode,
			strings.Replace(extendedTemplateCode, "<ENTITY>", kindName, -1),
		)
	}

	return extendedCode, nil
}
