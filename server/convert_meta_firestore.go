package main

import (
	"fmt"

	"cloud.google.com/go/firestore"
)

func convertMetaKeyToDocumentRef(
	client *firestore.Client,
	key *Key,
) (*firestore.DocumentRef, error) {
	if key == nil || key.PartitionId == nil {
		return nil, fmt.Errorf("key or key partition ID is nil; if the caller wants to allow nil keys, it must check to see if the input key is nil first")
	}

	namespace := key.PartitionId.Namespace

	firestoreTestCollection := client.Collection("Test")
	firestoreNamespace := firestoreTestCollection.Path[0:(len(firestoreTestCollection.Path) - len(firestoreTestCollection.ID) - 1)]

	if namespace == "" {
		namespace = firestoreNamespace
	}
	if namespace != firestoreNamespace {
		return nil, fmt.Errorf("namespace must be either omitted, or match '%s' for this Firestore-backed entity", firestoreNamespace)
	}

	var ref *firestore.DocumentRef
	for _, pathElement := range key.Path {
		if ref == nil {
			ref = client.Collection(pathElement.Kind).
				Doc(pathElement.GetName())
		} else {
			ref = ref.Collection(pathElement.Kind).
				Doc(pathElement.GetName())
		}
	}
	if ref == nil {
		return nil, fmt.Errorf("inbound key did not contain any path components: namespace '%s'", namespace)
	}

	return ref, nil
}
