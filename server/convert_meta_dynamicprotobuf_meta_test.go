package main

import (
	firebase "firebase.google.com/go"
	"github.com/jhump/protoreflect/dynamic"

	"context"

	"testing"

	"gotest.tools/assert"
)

func TestConvertMetaKeyToDynamicKeyToMetaKey(t *testing.T) {
	verifyKeyIntact(
		t,
		&Key{
			PartitionId: &PartitionId{
				Namespace: "mynamespace",
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
	)
}

func verifyKeyIntact(t *testing.T, key *Key) {
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

	// convert key
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
	assert.Assert(t, key.PartitionId != nil)
	assert.Assert(t, resultKey.PartitionId != nil)
	assert.Equal(t, key.PartitionId.Namespace, resultKey.PartitionId.Namespace)
	assert.Equal(t, len(key.Path), len(resultKey.Path))
}
