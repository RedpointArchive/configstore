package main

import (
	"cloud.google.com/go/firestore"
)

type operationProcessor struct {
	client *firestore.Client
	tx     *firestore.Transaction
}

func createOperationProcessor(client *firestore.Client, tx *firestore.Transaction) *operationProcessor {
	return &operationProcessor{
		client: client,
		tx:     tx,
	}
}
