package main

import (
	"cloud.google.com/go/firestore"
)

type operationProcessor struct {
	client *firestore.Client
}

func createOperationProcessor(client *firestore.Client) *operationProcessor {
	return &operationProcessor{
		client: client,
	}
}
