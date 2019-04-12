package main

import (
	firebase "firebase.google.com/go"
	"github.com/jhump/protoreflect/dynamic"

	"context"

	"testing"

	"gotest.tools/assert"
)

type KeyTestType int

const (
	KeyTestType_All           KeyTestType = 0
	KeyTestType_ProtobufOnly  KeyTestType = 1
	KeyTestType_FirestoreOnly KeyTestType = 2
)

func TestConvertKeyFullPathElement(t *testing.T) {
	verifyKeyIntact(
		t,
		&Key{
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
		KeyTestType_All,
	)
}

func TestConvertKeyUnsetFullPathElement(t *testing.T) {
	verifyKeyIntact(
		t,
		&Key{
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
				&PathElement{
					Kind:   "C",
					IdType: nil,
				},
			},
		},
		KeyTestType_ProtobufOnly,
	)
}

func TestConvertKeySinglePathElement(t *testing.T) {
	verifyKeyIntact(
		t,
		&Key{
			PartitionId: &PartitionId{
				Namespace: "projects/configstore-test-001/databases/(default)/documents",
			},
			Path: []*PathElement{
				&PathElement{
					Kind: "B",
					IdType: &PathElement_Name{
						Name: "MyName",
					},
				},
			},
		},
		KeyTestType_All,
	)
}

func TestConvertKeyUnsetSinglePathElement(t *testing.T) {
	verifyKeyIntact(
		t,
		&Key{
			PartitionId: &PartitionId{
				Namespace: "projects/configstore-test-001/databases/(default)/documents",
			},
			Path: []*PathElement{
				&PathElement{
					Kind:   "B",
					IdType: nil,
				},
			},
		},
		KeyTestType_ProtobufOnly,
	)
}

func TestConvertKeyNilPath(t *testing.T) {
	verifyKeyIntact(
		t,
		&Key{
			PartitionId: &PartitionId{
				Namespace: "projects/configstore-test-001/databases/(default)/documents",
			},
			Path: nil,
		},
		KeyTestType_ProtobufOnly,
	)
}

func verifyKeyIntact(t *testing.T, key *Key, ty KeyTestType) {
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

	if ty == KeyTestType_All || ty == KeyTestType_ProtobufOnly {
		// convert key meta -> dynamic -> meta
		messageFactory := dynamic.NewMessageFactoryWithDefaults()
		message, err := convertMetaKeyToDynamicKey(
			messageFactory,
			key,
			genResult.CommonMessageDescriptors,
		)
		assert.NilError(t, err)
		resultKey, err := convertDynamicKeyToMetaKey(
			client,
			message,
		)
		assert.NilError(t, err)

		// verify that it survived intact
		assert.DeepEqual(t, key, resultKey)
	}

	if ty == KeyTestType_All || ty == KeyTestType_FirestoreOnly {
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
	}
}
