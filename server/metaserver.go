package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

type configstoreMetaServiceServer struct {
	schema *Schema
}

func (s *configstoreMetaServiceServer) GetSchema(ctx context.Context, req *GetSchemaRequest) (*GetSchemaResponse, error) {
	return &GetSchemaResponse{
		Schema: s.schema,
	}, nil
}

func (s *configstoreMetaServiceServer) GetDefaultPartitionId(ctx context.Context, req *GetDefaultPartitionIdRequest) (*GetDefaultPartitionIdResponse, error) {
	firestoreTestCollection := client.Collection("Test")
	firestoreNamespace := firestoreTestCollection.Path[0:(len(firestoreTestCollection.Path) - len(firestoreTestCollection.ID) - 1)]

	return &GetDefaultPartitionIdResponse{
		Namespace: firestoreNamespace,
	}, nil
}

func convertSnapshotToMetaEntity(kindInfo *SchemaKind, snapshot *firestore.DocumentSnapshot) (*MetaEntity, error) {
	key, err := convertDocumentRefToMetaKey(snapshot.Ref)
	if err != nil {
		return nil, fmt.Errorf("error while converting firestore ref to meta key: %v", err)
	}
	entity := &MetaEntity{
		Key: key,
	}
	for key, value := range snapshot.Data() {
		for _, field := range kindInfo.Fields {
			if field.Name == key {
				f := &Value{
					Id: field.Id,
				}
				switch field.Type {
				case ValueType_double:
					switch v := value.(type) {
					case float64:
						f.DoubleValue = v
					default:
						f.DoubleValue = 0
					}
				case ValueType_int64:
					switch v := value.(type) {
					case int64:
						f.Int64Value = v
					default:
						f.Int64Value = 0
					}
				case ValueType_string:
					switch v := value.(type) {
					case string:
						f.StringValue = v
					default:
						f.StringValue = ""
					}
				case ValueType_timestamp:
					switch v := value.(type) {
					case []byte:
						f.TimestampValue = v
					default:
						f.TimestampValue = nil
					}
				case ValueType_boolean:
					switch v := value.(type) {
					case bool:
						f.BooleanValue = v
					default:
						f.BooleanValue = false
					}
				case ValueType_bytes:
					switch v := value.(type) {
					case []byte:
						f.BytesValue = v
					default:
						f.BytesValue = nil
					}
				case ValueType_key:
					switch v := value.(type) {
					case *firestore.DocumentRef:
						f.KeyValue, err = convertDocumentRefToMetaKey(v)
						if err != nil {
							f.KeyValue = nil
							fmt.Printf("error while converting firestore ref to meta key: %v", err)
						}
					default:
						f.KeyValue = nil
					}
				}
				entity.Values = append(entity.Values, f)
				break
			}
		}
	}
	return entity, nil
}

func (s *configstoreMetaServiceServer) MetaList(ctx context.Context, req *MetaListEntitiesRequest) (*MetaListEntitiesResponse, error) {
	var start interface{}
	if req.Start != nil {
		if len(req.Start[:]) > 0 {
			start = string(req.Start[:])
		}
	}

	var kindInfo *SchemaKind
	for kindName, kind := range s.schema.Kinds {
		if kindName == req.KindName {
			kindInfo = kind
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
		entity, err := convertSnapshotToMetaEntity(kindInfo, snapshot)
		if err != nil {
			fmt.Printf("%v", err)
			continue
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

func (s *configstoreMetaServiceServer) MetaGet(ctx context.Context, req *MetaGetEntityRequest) (*MetaGetEntityResponse, error) {
	var kindInfo *SchemaKind
	for kindName, kind := range s.schema.Kinds {
		if kindName == req.KindName {
			kindInfo = kind
			break
		}
	}
	if kindInfo == nil {
		return nil, fmt.Errorf("no such kind")
	}

	ref, err := convertMetaKeyToDocumentRef(
		client,
		req.Key,
	)
	if err != nil {
		return nil, err
	}

	snapshot, err := ref.Get(ctx)
	if err != nil {
		return nil, err
	}

	entity, err := convertSnapshotToMetaEntity(kindInfo, snapshot)
	if err != nil {
		return nil, err
	}

	response := &MetaGetEntityResponse{
		Entity: entity,
	}

	return response, nil
}

func (s *configstoreMetaServiceServer) MetaUpdate(ctx context.Context, req *MetaUpdateEntityRequest) (*MetaUpdateEntityResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *configstoreMetaServiceServer) MetaDelete(ctx context.Context, req *MetaDeleteEntityRequest) (*MetaDeleteEntityResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *configstoreMetaServiceServer) MetaCreate(ctx context.Context, req *MetaCreateEntityRequest) (*MetaCreateEntityResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
