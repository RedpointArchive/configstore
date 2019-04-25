package main

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

func convertMetaEntityToDynamicMessage(
	messageFactory *dynamic.MessageFactory,
	messageDescriptor *desc.MessageDescriptor,
	entity *MetaEntity,
	common *commonMessageDescriptors,
	schemaKind *SchemaKind,
) (*dynamic.Message, error) {
	out := messageFactory.NewDynamicMessage(messageDescriptor)
	out.SetFieldByName("key", entity.Key)

	for _, value := range entity.Values {
		field := findSchemaFieldByID(schemaKind, value.Id)
		if field == nil {
			// extra data not specified in the schema any more
			// we can safely ignore this
			continue
		}

		var err error
		if value == nil {
			err = out.TryClearFieldByName(field.Name)
		} else {
			switch field.Type {
			case ValueType_double:
				err = out.TrySetFieldByName(field.Name, value.DoubleValue)
				break
			case ValueType_int64:
				err = out.TrySetFieldByName(field.Name, value.Int64Value)
				break
			case ValueType_string:
				err = out.TrySetFieldByName(field.Name, value.StringValue)
				break
			case ValueType_timestamp:
				if value.TimestampValue == nil {
					err = out.TryClearFieldByName(field.Name)
				} else {
					err = out.TrySetFieldByName(field.Name, value.TimestampValue)
				}
				break
			case ValueType_boolean:
				err = out.TrySetFieldByName(field.Name, value.BooleanValue)
				break
			case ValueType_bytes:
				err = out.TrySetFieldByName(field.Name, value.BytesValue)
				break
			case ValueType_key:
				if value.KeyValue == nil {
					err = out.TryClearFieldByName(field.Name)
				} else {
					err = out.TrySetFieldByName(field.Name, value.KeyValue)
				}
				break
			case ValueType_uint64:
				err = out.TrySetFieldByName(field.Name, value.Uint64Value)
				break
			}
		}

		if err != nil {
			fmt.Printf("warning: encountered error while retrieving data from field '%s' on entity of kind '%s' with ID '%s' from Firestore: %v\n", field.Name, "todo", "todo", err)
		}
	}

	return out, nil
}
