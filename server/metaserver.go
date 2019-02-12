package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

type configstoreMetaServiceServer struct {
	schema *configstoreSchema
}

func convertType(t configstoreSchemaKindFieldType) ValueType {
	switch t {
	case typeDouble:
		return ValueType_Double
	case typeInt64:
		return ValueType_Int64
	case typeString:
		return ValueType_String
	case typeTimestamp:
		return ValueType_Timestamp
	case typeBoolean:
		return ValueType_Boolean
	}

	return ValueType_Double
}

func convertEditorType(t configstoreSchemaKindFieldEditorType) FieldEditorInfoType {
	switch t {
	case editorTypeDefault:
		return FieldEditorInfoType_Default
	case editorTypePassword:
		return FieldEditorInfoType_Password
	case editorTypeLookup:
		return FieldEditorInfoType_Lookup
	}

	return FieldEditorInfoType_Default
}

func (s *configstoreMetaServiceServer) GetSchema(ctx context.Context, req *GetSchemaRequest) (*GetSchemaResponse, error) {
	kinds := make([]*Kind, 0)
	for kindName, kind := range s.schema.Kinds {
		fields := make([]*Field, 0)
		for _, field := range kind.Fields {
			fields = append(fields, &Field{
				Id:      field.ID,
				Name:    field.Name,
				Type:    convertType(field.Type),
				Comment: field.Comment,
				Editor: &FieldEditorInfo{
					DisplayName: field.Editor.DisplayName,
					Type:        convertEditorType(field.Editor.Type),
					Readonly:    field.Editor.Readonly,
					ForeignType: field.Editor.ForeignType,
				},
			})
		}

		kinds = append(kinds, &Kind{
			Name:   kindName,
			Fields: fields,
			Editor: &KindEditor{
				Singular: kind.Editor.Singular,
				Plural:   kind.Editor.Plural,
			},
		})
	}

	schema := &Schema{
		Name:  s.schema.Name,
		Kinds: kinds,
	}

	return &GetSchemaResponse{
		Schema: schema,
	}, nil
}

func (s *configstoreMetaServiceServer) MetaList(ctx context.Context, req *MetaListEntitiesRequest) (*MetaListEntitiesResponse, error) {
	var start interface{}
	if req.Start != nil {
		if len(req.Start[:]) > 0 {
			start = string(req.Start[:])
		}
	}

	var kindInfo *configstoreSchemaKind
	for kindName, kind := range s.schema.Kinds {
		if kindName == req.KindName {
			kindInfo = &kind
			break
		}
	}
	if kindInfo == nil {
		return nil, fmt.Errorf("no such kind")
	}

	var err error
	var snapshots []*firestore.DocumentSnapshot
	if (req.Limit == 0) && start == nil {
		snapshots, err = client.Collection(req.KindName).Documents(ctx).GetAll()
	} else if req.Limit == 0 {
		snapshots, err = client.Collection(req.KindName).OrderBy(firestore.DocumentID, firestore.Asc).StartAfter(start.(string)).Documents(ctx).GetAll()
	} else if start == nil {
		snapshots, err = client.Collection(req.KindName).Limit(int(req.Limit)).Documents(ctx).GetAll()
	} else {
		snapshots, err = client.Collection(req.KindName).OrderBy(firestore.DocumentID, firestore.Asc).StartAfter(start.(string)).Limit(int(req.Limit)).Documents(ctx).GetAll()
	}

	if err != nil {
		return nil, err
	}

	var entities []*MetaEntity
	for _, snapshot := range snapshots {
		entity := &MetaEntity{
			Key: &Key{
				Val:   snapshot.Ref.ID,
				IsSet: true,
			},
		}
		for key, value := range snapshot.Data() {
			for _, field := range kindInfo.Fields {
				if field.Name == key {
					f := &Value{
						Id: field.ID,
					}
					switch field.Type {
					case typeDouble:
						switch v := value.(type) {
						case float64:
							f.DoubleValue = v
						default:
							f.DoubleValue = 0
						}
					case typeInt64:
						switch v := value.(type) {
						case int64:
							f.Int64Value = v
						default:
							f.Int64Value = 0
						}
					case typeString:
						switch v := value.(type) {
						case string:
							f.StringValue = v
						default:
							f.StringValue = ""
						}
					case typeTimestamp:
						switch v := value.(type) {
						case []byte:
							f.TimestampValue = v
						default:
							f.TimestampValue = nil
						}
					case typeBoolean:
						switch v := value.(type) {
						case bool:
							f.BooleanValue = v
						default:
							f.BooleanValue = false
						}
					case typeBytes:
						switch v := value.(type) {
						case []byte:
							f.BytesValue = v
						default:
							f.BytesValue = nil
						}
					case typeKey:
						switch v := value.(type) {
						case *firestore.DocumentRef:
							f.KeyValue = &Key{
								Val:   v.ID,
								IsSet: true,
							}
						default:
							f.KeyValue = &Key{
								Val:   "",
								IsSet: false,
							}
						}
					}
					entity.Values = append(entity.Values, f)
					break
				}
			}
		}
		entities = append(entities, entity)
	}

	response := &MetaListEntitiesResponse{
		Entities: entities,
	}

	if !(req.Limit == 0) {
		if uint32(len(entities)) < req.Limit {
			response.MoreResults = false
		} else {
			// TODO: query to see if there really are more results, to make this behave like datastore
			response.MoreResults = true
			last := snapshots[len(snapshots)-1]
			response.Next = []byte(last.Ref.ID)
		}
	} else {
		response.MoreResults = false
	}

	return response, nil
}
