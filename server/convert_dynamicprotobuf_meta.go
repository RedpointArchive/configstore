package main

import (
	"fmt"

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
			return nil, fmt.Errorf("unable to retrieve kind value from path element at index", idx)
		}
		nameField, _ := elem.TryGetFieldByName("name")
		idField, _ := elem.TryGetFieldByName("id")

		pathElement := &PathElement{
			Kind: kindField.(string),
		}

		if nameField != nil {
			pathElement.IdType = &PathElement_Name{
				Name: nameField.(string),
			}
		} else if idField != nil {
			pathElement.IdType = &PathElement_Id{
				Id: idField.(int64),
			}
		}

		pathElements = append(pathElements, pathElement)
	}

	if len(pathElements) == 0 {
		return nil, fmt.Errorf("key contained no path elements")
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
			metaEntity.Values = append(
				metaEntity.Values,
				&Value{
					Id:          fieldDescriptor.GetNumber(),
					DoubleValue: value,
				},
			)
			break
		case int64:
			metaEntity.Values = append(
				metaEntity.Values,
				&Value{
					Id:         fieldDescriptor.GetNumber(),
					Int64Value: value,
				},
			)
			break
		case string:
			metaEntity.Values = append(
				metaEntity.Values,
				&Value{
					Id:          fieldDescriptor.GetNumber(),
					StringValue: value,
				},
			)
			break
		case bool:
			metaEntity.Values = append(
				metaEntity.Values,
				&Value{
					Id:           fieldDescriptor.GetNumber(),
					BooleanValue: value,
				},
			)
			break
		case []byte:
			metaEntity.Values = append(
				metaEntity.Values,
				&Value{
					Id:         fieldDescriptor.GetNumber(),
					BytesValue: value,
				},
			)
			break
		case uint64:
			// We store uint64 as int64 inside Firestore, as Firestore
			// does not support uint64 natively
			metaEntity.Values = append(
				metaEntity.Values,
				&Value{
					Id:          fieldDescriptor.GetNumber(),
					Uint64Value: value,
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
					metaEntity.Values = append(
						metaEntity.Values,
						&Value{
							Id:       fieldDescriptor.GetNumber(),
							KeyValue: nkey,
						},
					)
				} else {
					metaEntity.Values = append(
						metaEntity.Values,
						&Value{
							Id:       fieldDescriptor.GetNumber(),
							KeyValue: nil,
						},
					)
				}
				break
			case "Timestamp":
				seconds := value.GetFieldByName("seconds")
				nanos := value.GetFieldByName("nanos")

				metaEntity.Values = append(
					metaEntity.Values,
					&Value{
						Id: fieldDescriptor.GetNumber(),
						TimestampValue: &timestamp.Timestamp{
							Seconds: seconds.(int64),
							Nanos:   nanos.(int32),
						},
					},
				)
				break
			default:
				return nil, fmt.Errorf("field '%s' contained unknown protobuf message", fieldDescriptor.GetName())
			}
		default:
			return nil, fmt.Errorf("field '%s' contained unknown field type", fieldDescriptor.GetName())
		}
	}

	return metaEntity, nil
}
