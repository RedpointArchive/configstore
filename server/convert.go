package main

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"

	"cloud.google.com/go/firestore"
)

func getTopLevelParent(ref *firestore.DocumentRef) *firestore.CollectionRef {
	var lastCollection *firestore.CollectionRef
	for ref != nil {
		lastCollection = ref.Parent
		ref = lastCollection.Parent
	}
	return lastCollection
}

func convertDocumentRefToKey(
	messageFactory *dynamic.MessageFactory,
	ref *firestore.DocumentRef,
	common *commonMessageDescriptors,
) (*dynamic.Message, error) {
	lastCollection := getTopLevelParent(ref)
	if lastCollection == nil {
		return nil, fmt.Errorf("ref has no top level parent")
	}

	partitionID := messageFactory.NewDynamicMessage(common.PartitionId)
	partitionID.SetFieldByName("namespace", lastCollection.Path[0:(len(lastCollection.Path)-len(lastCollection.ID)-1)])

	var reversePaths []*dynamic.Message
	for ref != nil {
		pathElement := messageFactory.NewDynamicMessage(common.PathElement)
		pathElement.SetFieldByName("kind", ref.Parent.ID)
		pathElement.SetFieldByName("name", ref.ID)

		reversePaths = append(reversePaths, pathElement)
	}

	var paths []*dynamic.Message
	for i := len(reversePaths) - 1; i >= 0; i-- {
		paths = append(paths, reversePaths[i])
	}

	key := messageFactory.NewDynamicMessage(common.Key)
	key.SetFieldByName("partitionId", partitionID)
	key.SetFieldByName("path", paths)

	return key, nil
}

func convertKeyToDocumentRef(
	client *firestore.Client,
	key *dynamic.Message,
) (*firestore.DocumentRef, error) {
	partitionID := key.GetFieldByName("partitionId")
	namespaceRaw := partitionID.(*dynamic.Message).GetFieldByName("namespace")
	namespace := namespaceRaw.(string)

	if namespace != "" {
		return nil, fmt.Errorf("namespace must be nil for Firestore-backed entity")
	}

	pathsRaw := key.GetFieldByName("paths")
	paths := pathsRaw.([]*dynamic.Message)

	var ref *firestore.DocumentRef
	for _, pathElement := range paths {
		if ref == nil {
			ref = client.Collection(pathElement.GetFieldByName("kind").(string)).
				Doc(pathElement.GetFieldByName("name").(string))
		} else {
			ref = ref.Collection(pathElement.GetFieldByName("kind").(string)).
				Doc(pathElement.GetFieldByName("name").(string))
		}
	}

	return ref, nil
}

func convertSnapshotToDynamicMessage(
	messageFactory *dynamic.MessageFactory,
	messageDescriptor *desc.MessageDescriptor,
	snapshot *firestore.DocumentSnapshot,
	common *commonMessageDescriptors,
) (*dynamic.Message, error) {
	key, err := convertDocumentRefToKey(
		messageFactory,
		snapshot.Ref,
		common,
	)
	if err != nil {
		return nil, err
	}

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

func convertDynamicMessageIntoKeyAndDataMap(
	client *firestore.Client,
	messageFactory *dynamic.MessageFactory,
	messageDescriptor *desc.MessageDescriptor,
	message *dynamic.Message,
	keyMessageDescriptor *desc.MessageDescriptor,
) (*firestore.DocumentRef, map[string]interface{}, error) {
	keyRaw, err := message.TryGetFieldByName("key")
	if err != nil {
		return nil, nil, err
	}

	key, err := convertKeyToDocumentRef(
		client,
		keyRaw,
	)
	if err != nil {
		return nil, nil, err
	}

	m := make(map[string]interface{})

	for _, fieldDescriptor := range message.GetKnownFields() {
		if fieldDescriptor.GetName() == "key" {
			continue
		}
		field := message.GetField(fieldDescriptor)
		m[fieldDescriptor.GetName()] = field
	}

	return key, m, nil
}
