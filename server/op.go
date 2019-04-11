package main

import (
	"cloud.google.com/go/firestore"
)

type operationProcessor struct {
	client *firestore.Client
}
