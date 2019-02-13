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

		ref = ref.Parent.Parent
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

	firestoreTestCollection := client.Collection("Test")
	firestoreNamespace := firestoreTestCollection.Path[0:(len(firestoreTestCollection.Path) - len(firestoreTestCollection.ID) - 1)]

	if namespace == "" {
		namespace = firestoreNamespace
	}
	if namespace != firestoreNamespace {
		return nil, fmt.Errorf("namespace must be either omitted, or match '%s' for this Firestore-backed entity", firestoreNamespace)
	}

	pathsRaw := key.GetFieldByName("path")
	pathsArray, ok := pathsRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("key path is not expected array of path elements")
	}
	var paths []*dynamic.Message
	for idx, e := range pathsArray {
		pe, ok := e.(*dynamic.Message)
		if !ok {
			return nil, fmt.Errorf("key path is not expected array of path elements (element %d)", idx)
		} else {
			paths = append(paths, pe)
		}
	}

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

func convertDocumentRefToMetaKey(
	ref *firestore.DocumentRef,
) (*Key, error) {
	lastCollection := getTopLevelParent(ref)
	if lastCollection == nil {
		return nil, fmt.Errorf("ref has no top level parent")
	}

	partitionID := &PartitionId{
		Namespace: lastCollection.Path[0:(len(lastCollection.Path) - len(lastCollection.ID) - 1)],
	}

	var reversePaths []*PathElement
	for ref != nil {
		pathElement := &PathElement{
			Kind: ref.Parent.ID,
			IdType: &PathElement_Name{
				Name: ref.ID,
			},
		}

		reversePaths = append(reversePaths, pathElement)
	}

	var paths []*PathElement
	for i := len(reversePaths) - 1; i >= 0; i-- {
		paths = append(paths, reversePaths[i])
	}

	return &Key{
		PartitionId: partitionID,
		Path:        paths,
	}, nil
}

func convertMetaKeyToDocumentRef(
	client *firestore.Client,
	key *Key,
) (*firestore.DocumentRef, error) {
	return nil, fmt.Errorf("not implemented")
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

func convertDynamicMessageIntoRefAndDataMap(
	client *firestore.Client,
	messageFactory *dynamic.MessageFactory,
	messageDescriptor *desc.MessageDescriptor,
	message *dynamic.Message,
) (*firestore.DocumentRef, map[string]interface{}, error) {
	keyRaw, err := message.TryGetFieldByName("key")
	if err != nil {
		return nil, nil, err
	}

	keyCon, ok := keyRaw.(*dynamic.Message)
	if !ok {
		return nil, nil, fmt.Errorf("key of unexpected type")
	}

	key, err := convertKeyToDocumentRef(
		client,
		keyCon,
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
