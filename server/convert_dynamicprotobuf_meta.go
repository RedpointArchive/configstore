package main

import (
	"fmt"
	"sort"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"

	"cloud.google.com/go/firestore"
)

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

	key, ok := keyRaw.(*Key)
	if !ok {
		return nil, fmt.Errorf("key of unexpected type")
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
		case *Key:
			setFields[fieldDescriptor.GetNumber()] = true
			metaEntity.Values = append(
				metaEntity.Values,
				&Value{
					Id:       fieldDescriptor.GetNumber(),
					Type:     ValueType_key,
					KeyValue: value,
				},
			)
			break
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
