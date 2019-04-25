package main

import (
	"fmt"
	"strings"

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

func serializeRef(ref *firestore.DocumentRef) string {
	key, err := convertDocumentRefToMetaKey(ref)
	if err != nil {
		return ""
	}
	return serializeKey(key)
}

func serializeKey(key *Key) string {
	if key == nil {
		return ""
	}

	var elements []string
	for _, pathElement := range key.Path {
		if _, ok := pathElement.IdType.(*PathElement_Id); ok {
			elements = append(elements, fmt.Sprintf("%s:id=%d", pathElement.GetKind(), pathElement.GetId()))
		} else if _, ok := pathElement.IdType.(*PathElement_Name); ok {
			elements = append(elements, fmt.Sprintf("%s:name=%s", pathElement.GetKind(), pathElement.GetName()))
		} else {
			elements = append(elements, fmt.Sprintf("%s:unset", pathElement.GetKind()))
		}
	}
	return fmt.Sprintf("ns=%s|%s", key.PartitionId.Namespace, strings.Join(elements, "|"))
}
