package main

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"

	"cloud.google.com/go/firestore"
)

func convertSnapshotToDynamicMessage(
	messageFactory *dynamic.MessageFactory,
	messageDescriptor *desc.MessageDescriptor,
	snapshot *firestore.DocumentSnapshot,
) (*dynamic.Message, error) {
	out := messageFactory.NewDynamicMessage(messageDescriptor)

	for name, value := range snapshot.Data() {
		out.SetFieldByName(name, value)
	}

	return out, nil
}
