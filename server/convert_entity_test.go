package main

import (
	firebase "firebase.google.com/go"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jhump/protoreflect/dynamic"

	"context"

	"testing"

	"gotest.tools/assert"
)

type EntityTestType int

const (
	EntityTestType_All           EntityTestType = 0
	EntityTestType_ProtobufOnly  EntityTestType = 1
	EntityTestType_FirestoreOnly EntityTestType = 2
)

func TestConvertEntity(t *testing.T) {
	verifyEntityIntact(
		t,
		&MetaEntity{
			Key: &Key{
				PartitionId: &PartitionId{
					Namespace: "projects/configstore-test-001/databases/(default)/documents",
				},
				Path: []*PathElement{
					&PathElement{
						Kind: "A",
						IdType: &PathElement_Id{
							Id: 9223372036854775807,
						},
					},
					&PathElement{
						Kind: "B",
						IdType: &PathElement_Name{
							Name: "MyName",
						},
					},
				},
			},
			Values: []*Value{
				&Value{
					Id:          2,
					Type:        ValueType_double,
					DoubleValue: 456,
				},
				&Value{
					Id:         3,
					Type:       ValueType_int64,
					Int64Value: 123,
				},
				&Value{
					Id:          4,
					Type:        ValueType_string,
					StringValue: "hello world",
				},
				&Value{
					Id:   5,
					Type: ValueType_timestamp,
					TimestampValue: &timestamp.Timestamp{
						Seconds: 123,
						Nanos:   456,
					},
				},
				&Value{
					Id:           6,
					Type:         ValueType_boolean,
					BooleanValue: true,
				},
				&Value{
					Id:         7,
					Type:       ValueType_bytes,
					BytesValue: []byte("some bytes"),
				},
				&Value{
					Id:   8,
					Type: ValueType_key,
					KeyValue: &Key{
						PartitionId: &PartitionId{
							Namespace: "projects/configstore-test-001/databases/(default)/documents",
						},
						Path: []*PathElement{
							&PathElement{
								Kind: "C",
								IdType: &PathElement_Id{
									Id: 9213372036854775807,
								},
							},
							&PathElement{
								Kind: "D",
								IdType: &PathElement_Name{
									Name: "AnotherName",
								},
							},
						},
					},
				},
				&Value{
					Id:       9,
					Type:     ValueType_key,
					KeyValue: nil,
				},
				&Value{
					Id:         10,
					Type:       ValueType_int64,
					Int64Value: -123,
				},
				&Value{
					Id:          11,
					Type:        ValueType_uint64,
					Uint64Value: 123,
				},
				&Value{
					Id:             12,
					Type:           ValueType_timestamp,
					TimestampValue: nil,
				},
			},
		},
		EntityTestType_All,
	)
}

func verifyEntityIntact(t *testing.T, entity *MetaEntity, ty EntityTestType) {
	ctx := context.Background()

	// connect to Firestore
	conf := &firebase.Config{ProjectID: "configstore-test-001"}
	app, err := firebase.NewApp(ctx, conf)
	assert.NilError(t, err)
	client, err := app.Firestore(ctx)
	assert.NilError(t, err)
	defer client.Close()

	// generate schema
	genResult, err := generate("./schema.json")
	assert.NilError(t, err)

	if ty == EntityTestType_All || ty == EntityTestType_ProtobufOnly {
		// convert entity meta -> dynamic -> meta
		messageFactory := dynamic.NewMessageFactoryWithDefaults()
		message, err := convertMetaEntityToDynamicMessage(
			messageFactory,
			genResult.MessageMap["UnitTest001"],
			entity,
			genResult.CommonMessageDescriptors,
			genResult.KindMap[genResult.ServiceMap["UnitTest001"]],
		)
		assert.NilError(t, err)
		resultEntity, err := convertDynamicMessageIntoMetaEntity(
			client,
			messageFactory,
			genResult.MessageMap["UnitTest001"],
			message,
			genResult.KindMap[genResult.ServiceMap["UnitTest001"]],
		)
		assert.NilError(t, err)

		// verify that it survived intact
		assert.DeepEqual(t, entity, resultEntity)
	}
	/*
		if ty == All || ty == FirestoreOnly {
			// convert key meta -> firestore -> meta
			ref, err := convertMetaKeyToDocumentRef(
				client,
				key,
			)
			assert.NilError(t, err)
			resultKey2, err := convertDocumentRefToMetaKey(
				ref,
			)
			assert.NilError(t, err)

			// verify that it survived intact
			assert.DeepEqual(t, key, resultKey2)
		}*/
}
