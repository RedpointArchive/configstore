package main

import (
	"cloud.google.com/go/firestore"
)

type transactionProcessor struct {
	client *firestore.Client
}

func createTransactionProcessor(client *firestore.Client) *transactionProcessor {
	return &transactionProcessor{
		client: client,
	}
}
