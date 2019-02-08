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
	out.SetFieldByName("id", snapshot.Ref.ID)

	for name, value := range snapshot.Data() {
		out.SetFieldByName(name, value)
	}

	return out, nil
}

func convertDynamicMessageIntoIDAndDataMap(
	messageFactory *dynamic.MessageFactory,
	messageDescriptor *desc.MessageDescriptor,
	message *dynamic.Message,
) (string, map[string]interface{}, error) {
	idRaw, err := message.TryGetFieldByName("id")
	if err != nil {
		return "", nil, err
	}

	id := idRaw.(string)

	m := make(map[string]interface{})

	for _, fieldDescriptor := range message.GetKnownFields() {
		if fieldDescriptor.GetName() == "id" {
			continue
		}
		field := message.GetField(fieldDescriptor)
		m[fieldDescriptor.GetName()] = field
	}

	return id, m, nil
}
