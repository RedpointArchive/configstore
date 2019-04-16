package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

type transactionWatcher struct {
	client *firestore.Client
	schema *Schema

	currentEntities               map[string]*firestore.DocumentSnapshot
	currentEntitiesLock           sync.RWMutex
	pendingChangesByTimestamp     map[string][]*firestore.DocumentChange
	pendingChangesByTimestampLock sync.RWMutex
	inboundChanges                chan firestore.DocumentChange

	transactions     []*MetaTransactionRecord
	transactionsLock sync.RWMutex
}

func serializeTime(ts time.Time) string {
	return fmt.Sprintf("%d-%d", ts.Unix(), ts.Nanosecond())
}

func createTransactionWatcher(ctx context.Context, client *firestore.Client, schema *Schema) (*transactionWatcher, error) {
	watcher := &transactionWatcher{
		client:                    client,
		schema:                    schema,
		currentEntities:           make(map[string]*firestore.DocumentSnapshot),
		pendingChangesByTimestamp: make(map[string][]*firestore.DocumentChange),
		inboundChanges:            make(chan firestore.DocumentChange),
		transactions:              nil,
	}

	watcher.currentEntitiesLock.Lock()
	defer watcher.currentEntitiesLock.Unlock()

	// process transactions every second
	go func() {
		for true {
			// wait a second
			time.Sleep(time.Second * 1)

			// preemptive lock to check if we have any transactions to process at all
			watcher.transactionsLock.RLock()
			if len(watcher.transactions) == 0 {
				watcher.transactionsLock.RUnlock()
				continue
			}
			watcher.transactionsLock.RUnlock()

			// obtain all locks
			watcher.transactionsLock.Lock()
			watcher.pendingChangesByTimestampLock.Lock()
			watcher.currentEntitiesLock.Lock()
			defer watcher.transactionsLock.Unlock()
			defer watcher.pendingChangesByTimestampLock.Unlock()
			defer watcher.currentEntitiesLock.Unlock()

			// we have transactions to process
			for _, transaction := range watcher.transactions {
				// WIP: todo complete....
			}
		}
	}

	// listen for inbound changes
	go func() {
		for elem := range watcher.inboundChanges {
			watcher.pendingChangesByTimestampLock.Lock()
			defer watcher.pendingChangesByTimestampLock.Unlock()

			ts := serializeTime(elem.Doc.UpdateTime)

			watcher.pendingChangesByTimestamp[ts] = append(
				watcher.pendingChangesByTimestamp[ts],
				&elem,
			)
		}
	}()

	// listen for new transactions coming in
	go func() {
		transactions := watcher.client.Collection("Transaction").Snapshots(ctx)
		for true {
			transactionSnapshot, err := transactions.Next()
			if err != nil {
				log.Printf("error during transaction watch: %v", err)
				return
			}

			watcher.transactionsLock.Lock()
			defer watcher.transactionsLock.Unlock()

			for _, change := range transactionSnapshot.Changes {
				data := change.Data()

				var mutatedRefs []*firestore.DocumentRef
				var mutatedKeys []*Key
				var deletedRefs []*firestore.DocumentRef
				var deletedKeys []*Key
				var dateSubmitted time.Time
				var dateCreated time.Time
				var description string
				if w, ok := data["mutatedKeys"].([]*firestore.DocumentRef); ok {
					mutatedRefs = w
				}
				if w, ok := data["deletedKeys"].([]*firestore.DocumentRef); ok {
					deletedRefs = w
				}
				if w, ok := data["dateSubmitted"].(time.Time); ok {
					dateSubmitted = w
				}
				if w, ok := data["dateCreated"].(time.Time); ok {
					dateCreated = w
				}
				if w, ok := data["description"].(string); ok {
					description = w
				}
				for _, ref := mutatedRefs {
					key, err := convertDocumentRefToMetaKey(ref)
					if err != nil {
						continue
					}
					mutatedKeys = append(mutatedKeys, key)
				}
				for _, ref := deletedRefs {
					key, err := convertDocumentRefToMetaKey(ref)
					if err != nil {
						continue
					}
					deletedKeys = append(deletedKeys, key)
				}

				watcher.transactions = append(
					watcher.transactions,
					&MetaTransactionRecord{
						MutatedKeys: mutatedKeys,
						DeletedKeys: deletedKeys,
						DateSubmitted: convertTimeToTimestamp(dateSubmitted),
						DateCreated: convertTimeToTimestamp(dateCreated),
						Description: description,
					},
				)
			}

			sort.Slice(watcher.transactions[:], func(i, j int) bool {
				if (watcher.transactions[i].DateCreated.Seconds < watcher.transactions[j].DateCreated.Seconds) {
					return true
				}
				if (watcher.transactions[i].DateCreated.Nanos < watcher.transactions[j].DateCreated.Nanos) {
					return true
				}
				return false
			})
		}
	}()

	// for each kind, start watching the collection and pipe snapshots into
	// the inboundChanges channel
	for kindName := range schema.Kinds {
		snapshots := watcher.client.Collection(kindName).Snapshots(ctx)
		go func() {
			for true {
				snapshot, err := snapshots.Next()
				if err != nil {
					log.Printf("error during entity watch: %v", err)
					return
				}

				for _, change := range snapshot.Changes {
					watcher.inboundChanges <- change
				}
			}
		}()
	}

	// for each kind, fill in the current entities
	for kindName := range schema.Kinds {
		documents, err := watcher.client.Collection(kindName).Documents(ctx).GetAll()
		if err != nil {
			return nil, err
		}
		for _, document := range documents {
			watcher.currentEntities[serializeRef(document.Ref)] = document
		}
	}

	return watcher, nil
}
