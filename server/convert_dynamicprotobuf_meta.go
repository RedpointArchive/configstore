package main

import (
	"fmt"
	"sort"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"

	"cloud.google.com/go/firestore"
)

func convertDynamicKeyToMetaKey(
	client *firestore.Client,
	key *dynamic.Message,
) (*Key, error) {
	partitionID := key.GetFieldByName("partitionId")
	namespaceRaw := partitionID.(*dynamic.Message).GetFieldByName("namespace")
	namespace := namespaceRaw.(string)

	firestoreTestCollection := client.Collection("Test")
	firestoreNamespace := firestoreTestCollection.Path[0:(len(firestoreTestCollection.Path) - len(firestoreTestCollection.ID) - 1)]

	if namespace == "" {
		namespace = firestoreNamespace
	}
	if namespace != firestoreNamespace {
		return nil, fmt.Errorf("namespace must be either omitted, or match '%s' for this Firestore-backed entity", firestoreNamespace)
	}

	pathsRaw := key.GetFieldByName("path")
	pathsArray, ok := pathsRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("key path is not expected array of path elements")
	}
	var paths []*dynamic.Message
	for idx, e := range pathsArray {
		pe, ok := e.(*dynamic.Message)
		if !ok {
			return nil, fmt.Errorf("key path is not expected array of path elements (element %d)", idx)
		} else {
			paths = append(paths, pe)
		}
	}

	var pathElements []*PathElement
	for idx, elem := range paths {
		kindField, err := elem.TryGetFieldByName("kind")
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve kind value from path element at index %d", idx)
		}
		nameField, _ := elem.TryGetFieldByName("name")
		idField, _ := elem.TryGetFieldByName("id")

		pathElement := &PathElement{
			Kind: kindField.(string),
		}

		if elem.HasFieldName("name") {
			pathElement.IdType = &PathElement_Name{
				Name: nameField.(string),
			}
		} else if elem.HasFieldName("id") {
			pathElement.IdType = &PathElement_Id{
				Id: idField.(int64),
			}
		}

		pathElements = append(pathElements, pathElement)
	}

	return &Key{
		PartitionId: &PartitionId{
			Namespace: namespace,
		},
		Path: pathElements,
	}, nil
}

func convertDynamicMessageIntoMetaEntity(
	client *firestore.Client,
	messageFactory *dynamic.MessageFactory,
	messageDescriptor *desc.MessageDescriptor,
	message *dynamic.Message,
	schemaKind *SchemaKind,
) (*MetaEntity, error) {
	keyRaw, err := message.TryGetFieldByName("key")
	if err != nil {
		return nil, err
	}

	keyCon, ok := keyRaw.(*dynamic.Message)
	if !ok {
		return nil, fmt.Errorf("key of unexpected type")
	}

	key, err := convertDynamicKeyToMetaKey(
		client,
		keyCon,
	)
	if err != nil {
		return nil, err
	}

	metaEntity := &MetaEntity{
		Key:    key,
		Values: nil,
	}

	setFields := make(map[int32]bool)
	for _, fieldDescriptor := range message.GetKnownFields() {
		if fieldDescriptor.GetName() == "key" {
			continue
		}

		if !message.HasField(fieldDescriptor) {
			continue
		}

		rawValue := message.GetField(fieldDescriptor)

		switch value := rawValue.(type) {
		case float64:
			setFields[fieldDescriptor.GetNumber()] = true
			metaEntity.Values = append(
				metaEntity.Values,
				&Value{
					Id:          fieldDescriptor.GetNumber(),
					Type:        ValueType_double,
					DoubleValue: value,
				},
			)
			break
		case int64:
			setFields[fieldDescriptor.GetNumber()] = true
			metaEntity.Values = append(
				metaEntity.Values,
				&Value{
					Id:         fieldDescriptor.GetNumber(),
					Type:       ValueType_int64,
					Int64Value: value,
				},
			)
			break
		case string:
			setFields[fieldDescriptor.GetNumber()] = true
			metaEntity.Values = append(
				metaEntity.Values,
				&Value{
					Id:          fieldDescriptor.GetNumber(),
					Type:        ValueType_string,
					StringValue: value,
				},
			)
			break
		case bool:
			setFields[fieldDescriptor.GetNumber()] = true
			metaEntity.Values = append(
				metaEntity.Values,
				&Value{
					Id:           fieldDescriptor.GetNumber(),
					Type:         ValueType_boolean,
					BooleanValue: value,
				},
			)
			break
		case []byte:
			setFields[fieldDescriptor.GetNumber()] = true
			metaEntity.Values = append(
				metaEntity.Values,
				&Value{
					Id:         fieldDescriptor.GetNumber(),
					Type:       ValueType_bytes,
					BytesValue: value,
				},
			)
			break
		case uint64:
			// We store uint64 as int64 inside Firestore, as Firestore
			// does not support uint64 natively
			setFields[fieldDescriptor.GetNumber()] = true
			metaEntity.Values = append(
				metaEntity.Values,
				&Value{
					Id:          fieldDescriptor.GetNumber(),
					Type:        ValueType_uint64,
					Uint64Value: value,
				},
			)
			break
		case *timestamp.Timestamp:
			setFields[fieldDescriptor.GetNumber()] = true
			metaEntity.Values = append(
				metaEntity.Values,
				&Value{
					Id:             fieldDescriptor.GetNumber(),
					Type:           ValueType_timestamp,
					TimestampValue: value,
				},
			)
			break
		case *dynamic.Message:
			switch value.GetMessageDescriptor().GetName() {
			case "Key":
				partitionIDFd := value.GetMessageDescriptor().FindFieldByName("partitionId")
				pathFd := value.GetMessageDescriptor().FindFieldByName("path")

				if value.HasField(partitionIDFd) && value.HasField(pathFd) {
					nkey, err := convertDynamicKeyToMetaKey(
						client,
						value,
					)
					if err != nil {
						return nil, fmt.Errorf("error on field '%s': %v", fieldDescriptor.GetName(), err)
					}
					setFields[fieldDescriptor.GetNumber()] = true
					metaEntity.Values = append(
						metaEntity.Values,
						&Value{
							Id:       fieldDescriptor.GetNumber(),
							Type:     ValueType_key,
							KeyValue: nkey,
						},
					)
				} else {
					setFields[fieldDescriptor.GetNumber()] = true
					metaEntity.Values = append(
						metaEntity.Values,
						&Value{
							Id:       fieldDescriptor.GetNumber(),
							Type:     ValueType_key,
							KeyValue: nil,
						},
					)
				}
				break
			default:
				return nil, fmt.Errorf("field '%s' contained unknown protobuf message", fieldDescriptor.GetName())
			}
		default:
			return nil, fmt.Errorf("field '%s' contained unknown field type '%T' with value: %v", fieldDescriptor.GetName(), rawValue, rawValue)
		}
	}

	// polyfill nil fields (keys and timestamps)
	for _, schemaField := range schemaKind.Fields {
		if _, ok := setFields[schemaField.Id]; !ok {
			// need to polyfill this value
			if schemaField.Type == ValueType_key {
				metaEntity.Values = append(
					metaEntity.Values,
					&Value{
						Id:       schemaField.Id,
						Type:     ValueType_key,
						KeyValue: nil,
					},
				)
			} else if schemaField.Type == ValueType_timestamp {
				metaEntity.Values = append(
					metaEntity.Values,
					&Value{
						Id:             schemaField.Id,
						Type:           ValueType_timestamp,
						TimestampValue: nil,
					},
				)
			}
		}
	}

	// sort, this is mainly so unit tests pass...
	sort.Slice(metaEntity.Values[:], func(i, j int) bool {
		return metaEntity.Values[i].Id < metaEntity.Values[j].Id
	})

	return metaEntity, nil
}
