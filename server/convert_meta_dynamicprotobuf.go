package main

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

func convertMetaKeyToDynamicKey(
	messageFactory *dynamic.MessageFactory,
	key *Key,
	common *commonMessageDescriptors,
) (*dynamic.Message, error) {
	if key == nil {
		return nil, fmt.Errorf("meta key was nil, internal caller must handle this scenario\n")
	}

	partitionID := messageFactory.NewDynamicMessage(common.PartitionId)
	partitionID.SetFieldByName("namespace", key.PartitionId.Namespace)

	var path []*dynamic.Message
	for _, elem := range key.Path {
		pathElement := messageFactory.NewDynamicMessage(common.PathElement)
		pathElement.SetFieldByName("kind", elem.Kind)
		switch elem.GetIdType().(type) {
		case *PathElement_Id:
			pathElement.SetFieldByName("id", elem.GetId())
			break
		case *PathElement_Name:
			pathElement.SetFieldByName("name", elem.GetName())
			break
		}

		path = append(path, pathElement)
	}

	dynamicKey := messageFactory.NewDynamicMessage(common.Key)
	dynamicKey.SetFieldByName("partitionId", partitionID)
	dynamicKey.SetFieldByName("path", path)

	return dynamicKey, nil
}

func convertMetaEntityToDynamicMessage(
	messageFactory *dynamic.MessageFactory,
	messageDescriptor *desc.MessageDescriptor,
	entity *MetaEntity,
	common *commonMessageDescriptors,
	schemaKind *SchemaKind,
) (*dynamic.Message, error) {
	key, err := convertMetaKeyToDynamicKey(
		messageFactory,
		entity.Key,
		common,
	)
	if err != nil {
		return nil, err
	}

	out := messageFactory.NewDynamicMessage(messageDescriptor)
	out.SetFieldByName("key", key)

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
					key, err := convertMetaKeyToDynamicKey(
						messageFactory,
						value.KeyValue,
						common,
					)
					if err != nil {
						// pass error through
					} else {
						err = out.TrySetFieldByName(field.Name, key)
					}
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
