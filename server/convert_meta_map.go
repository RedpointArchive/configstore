package main

import (
	"cloud.google.com/go/firestore"
)

func convertMetaEntityToRefAndDataMap(
	client *firestore.Client,
	entity *MetaEntity,
	schema *SchemaKind,
) (*firestore.DocumentRef, map[string]interface{}, error) {
	key, err := convertMetaKeyToDocumentRef(
		client,
		entity.Key,
	)
	if err != nil {
		return nil, nil, err
	}

	m := make(map[string]interface{})

	for _, value := range entity.Values {
		name := ""
		for _, field := range schema.Fields {
			if field.Id == value.Id {
				name = field.Name
				break
			}
		}
		if name == "" {
			// field not found?
			continue
		}

		switch value.Type {
		case ValueType_double:
			m[name] = value.DoubleValue
			break
		case ValueType_int64:
			m[name] = value.Int64Value
			break
		case ValueType_string:
			m[name] = value.StringValue
			break
		case ValueType_timestamp:
			m[name] = value.TimestampValue
			break
		case ValueType_boolean:
			m[name] = value.BooleanValue
			break
		case ValueType_bytes:
			m[name] = value.BytesValue
			break
		case ValueType_key:
			if value.KeyValue == nil {
				m[name] = nil
			} else {
				ref, err := convertMetaKeyToDocumentRef(
					client,
					value.KeyValue,
				)
				if err != nil {
					return nil, nil, err
				}
				m[name] = ref
			}
			break
		case ValueType_uint64:
			// We store uint64 as int64 inside Firestore, as Firestore
			// does not support uint64 natively
			m[name] = int64(value.Uint64Value)
			break
		default:
			// todo, log?
			break
		}
	}

	// TODO: implement readonly safety
	/*
		if currentSnapshot != nil {
			for _, field := range schemaKind.Fields {
				if field.Readonly {
					// Verify that the property hasn't changed.
					if reflect.DeepEqual(m[field.Name], currentSnapshot.Data()[field.Name]) {
						return nil, nil, fmt.Errorf("readonly field '%s' contains mutated value", field.Name)
					}
				}
			}
		}
	*/

	return key, m, nil
}
