package main

import (
	"cloud.google.com/go/firestore"
)

type transactionProcessor struct {
	client *firestore.Client
}
