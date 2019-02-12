package main

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"

	"cloud.google.com/go/firestore"
)

func convertSnapshotToDynamicMessage(
	messageFactory *dynamic.MessageFactory,
	messageDescriptor *desc.MessageDescriptor,
	snapshot *firestore.DocumentSnapshot,
	keyMessageDescriptor *desc.MessageDescriptor,
) (*dynamic.Message, error) {
	key := messageFactory.NewDynamicMessage(keyMessageDescriptor)
	key.SetFieldByName("val", snapshot.Ref.ID)
	key.SetFieldByName("isSet", true)

	out := messageFactory.NewDynamicMessage(messageDescriptor)
	out.SetFieldByName("key", key)

	for name, value := range snapshot.Data() {
		fd := out.FindFieldDescriptorByName(name)
		if fd == nil {
			// extra data not specified in the schema any more
			// we can safely ignore this
		}

		err := out.TrySetFieldByName(name, value)
		if err != nil {
			fmt.Printf("warning: encountered error while retrieving data from field '%s' on entity of kind '%s' with ID '%s' from Firestore: %v\n", name, snapshot.Ref.Parent.ID, snapshot.Ref.ID, err)
		}
	}

	return out, nil
}

func convertDynamicMessageIntoID(
	messageFactory *dynamic.MessageFactory,
	keyMessageDescriptor *desc.MessageDescriptor,
	message *dynamic.Message,
) (string, error) {
	valRaw, err := message.TryGetFieldByName("val")
	if err != nil {
		return "", err
	}
	switch vv := valRaw.(type) {
	case string:
		return vv, nil
	default:
		return "", nil
	}
}

func convertDynamicMessageIntoIDAndDataMap(
	messageFactory *dynamic.MessageFactory,
	messageDescriptor *desc.MessageDescriptor,
	message *dynamic.Message,
	keyMessageDescriptor *desc.MessageDescriptor,
) (string, map[string]interface{}, error) {
	keyRaw, err := message.TryGetFieldByName("key")
	if err != nil {
		return "", nil, err
	}

	var id string
	switch v := keyRaw.(type) {
	case *dynamic.Message:
		valRaw, err := v.TryGetFieldByName("val")
		if err != nil {
			return "", nil, err
		}
		switch vv := valRaw.(type) {
		case string:
			id = vv
			break
		default:
			id = ""
			break
		}
		break
	default:
		id = ""
		break
	}

	m := make(map[string]interface{})

	for _, fieldDescriptor := range message.GetKnownFields() {
		if fieldDescriptor.GetName() == "key" {
			continue
		}
		field := message.GetField(fieldDescriptor)
		m[fieldDescriptor.GetName()] = field
	}

	return id, m, nil
}
