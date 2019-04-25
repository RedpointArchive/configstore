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
		var collectionRef *firestore.CollectionRef
		if ref == nil {
			collectionRef = client.Collection(pathElement.Kind)
		} else {
			collectionRef = ref.Collection(pathElement.Kind)
		}

		if pathElement.IdType == nil {
			// unset, automatically generate an ID
			ref = collectionRef.NewDoc()
		} else {
			switch pathElement.IdType.(type) {
			case *PathElement_Name:
				ref = collectionRef.Doc(pathElement.GetName())
				break
			case *PathElement_Id:
				ref = collectionRef.Doc(fmt.Sprintf("__datastore_id_polyfill=%d", pathElement.GetId()))
				break
			}
		}
	}
	if ref == nil {
		return nil, fmt.Errorf("inbound key did not contain any path components: namespace '%s'", namespace)
	}

	return ref, nil
}
