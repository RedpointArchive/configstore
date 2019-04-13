package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
)

func convertDocumentRefToMetaKey(
	ref *firestore.DocumentRef,
) (*Key, error) {
	lastCollection := getTopLevelParent(ref)
	if lastCollection == nil {
		return nil, fmt.Errorf("ref has no top level parent")
	}

	partitionID := &PartitionId{
		Namespace: lastCollection.Path[0:(len(lastCollection.Path) - len(lastCollection.ID) - 1)],
	}

	var reversePaths []*PathElement
	for ref != nil {
		var pathElement *PathElement
		if strings.HasPrefix(ref.ID, "__datastore_id_polyfill=") {
			id, _ := strconv.ParseInt(ref.ID[len("__datastore_id_polyfill="):], 10, 64)
			pathElement = &PathElement{
				Kind: ref.Parent.ID,
				IdType: &PathElement_Id{
					Id: id,
				},
			}
		} else {
			pathElement = &PathElement{
				Kind: ref.Parent.ID,
				IdType: &PathElement_Name{
					Name: ref.ID,
				},
			}
		}

		reversePaths = append(reversePaths, pathElement)
		ref = ref.Parent.Parent
	}

	var paths []*PathElement
	for i := len(reversePaths) - 1; i >= 0; i-- {
		paths = append(paths, reversePaths[i])
	}

	return &Key{
		PartitionId: partitionID,
		Path:        paths,
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
					Id:   field.Id,
					Type: field.Type,
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
				case ValueType_uint64:
					switch v := value.(type) {
					case int64:
						f.Uint64Value = uint64(v)
					default:
						f.Uint64Value = 0
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
					case time.Time:
						f.TimestampValue = convertTimeToTimestamp(v)
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
				default:
					return nil, fmt.Errorf("field type %d not supported in convertSnapshotToMetaEntity", field.Type)
				}
				entity.Values = append(entity.Values, f)
				break
			}
		}
	}

	// sort, this is mainly so unit tests pass...
	sort.Slice(entity.Values[:], func(i, j int) bool {
		return entity.Values[i].Id < entity.Values[j].Id
	})

	return entity, nil
}
