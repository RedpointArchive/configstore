package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
)

type transactionWatcher struct {
	client *firestore.Client
	schema *Schema

	currentEntities               map[string]*firestore.DocumentSnapshot
	currentEntitiesLock           sync.RWMutex
	pendingChangesByTimestamp     map[string]map[string]*firestore.DocumentChange
	pendingChangesByTimestampLock sync.RWMutex
	inboundChanges                chan firestore.DocumentChange

	transactions     []*MetaTransactionRecord
	transactionsLock sync.RWMutex

	outboundChanges chan *MetaTransactionBatch

	outboundChannels     []chan *MetaTransactionBatch
	outboundChannelsLock sync.Mutex
}

func (watcher *transactionWatcher) RegisterChannel(newCh chan *MetaTransactionBatch) {
	watcher.outboundChannelsLock.Lock()
	defer watcher.outboundChannelsLock.Unlock()

	for _, ch := range watcher.outboundChannels {
		if ch == newCh {
			return
		}
	}

	watcher.outboundChannels = append(
		watcher.outboundChannels,
		newCh,
	)
}

func (watcher *transactionWatcher) DeregisterChannel(oldCh chan *MetaTransactionBatch) {
	watcher.outboundChannelsLock.Lock()
	defer watcher.outboundChannelsLock.Unlock()

	idx := -1
	for i, ch := range watcher.outboundChannels {
		if ch == oldCh {
			idx = i
			break
		}
	}

	if idx == -1 {
		return
	}

	// remove the element quickly
	watcher.outboundChannels[len(watcher.outboundChannels)-1], watcher.outboundChannels[idx] = watcher.outboundChannels[idx], watcher.outboundChannels[len(watcher.outboundChannels)-1]
	watcher.outboundChannels = watcher.outboundChannels[:len(watcher.outboundChannels)-1]
}

func serializeTime(ts time.Time) string {
	return fmt.Sprintf("%d-%d", ts.Unix(), ts.Nanosecond())
}

func (watcher *transactionWatcher) processTransaction(i int, transaction *MetaTransactionRecord) bool {
	// check if the transaction is already complete
	allMutatedKeysAreNewer := true
	allDeletedKeysAreGoneOrCreatedAfterTransaction := true
	for _, mutatedKey := range transaction.MutatedKeys {
		mks := serializeKey(mutatedKey)
		mutatedEntity, ok := watcher.currentEntities[mks]
		if !ok {
			allMutatedKeysAreNewer = false
			break
		}
		if mutatedEntity.UpdateTime.Before(convertTimestampToTime(transaction.DateCreated)) {
			allMutatedKeysAreNewer = false
			break
		}
	}
	for _, deletedKey := range transaction.DeletedKeys {
		dks := serializeKey(deletedKey)
		deletedEntity, ok := watcher.currentEntities[dks]
		if ok {
			if deletedEntity.CreateTime.Before(convertTimestampToTime(transaction.DateCreated)) {
				allDeletedKeysAreGoneOrCreatedAfterTransaction = false
				break
			}
		}
	}
	if allMutatedKeysAreNewer && allDeletedKeysAreGoneOrCreatedAfterTransaction {
		// this transaction has already been applied to the current entities
		// list (usually this happens when we're loading initial state from
		// Firestore and have fetched the previous 5 minutes of transactions to
		// ensure coverage of transactions).
		return true
	}

	// at this point, we know the transaction hasn't already been applied to
	// current entities, so we need to check the pendingChangesByTimestamp
	// to see if we have all the mutated keys present (deleted keys won't exist;
	// we just remove them from current entities once we have all the mutations
	// we need)
	if len(transaction.MutatedKeys) > 0 {
		ts := serializeTime(convertTimestampToTime(transaction.DateCreated))
		pendingChanges, ok := watcher.pendingChangesByTimestamp[ts]
		if !ok {
			// we don't have any of the pending changes, but we have at least one
			// mutated key. we can't be ready to apply this transaction
			return false
		}
		for _, mutatedKey := range transaction.MutatedKeys {
			mks := serializeKey(mutatedKey)
			_, ok := pendingChanges[mks]
			if !ok {
				// we don't have this entity's snapshot yet
				return false
			}
		}
	}

	// at this point, we have all the changes for this transaction in pendingChanges,
	// we can apply them to currentEntities
	// TODO: emit the "transaction batch" data structure, containing all of the changes
	// for configstore clients to apply atomically
	batch := &MetaTransactionBatch{
		Description:     transaction.Description,
		DeletedKeys:     transaction.DeletedKeys,
		MutatedEntities: nil,
	}
	if len(transaction.MutatedKeys) > 0 {
		ts := serializeTime(convertTimestampToTime(transaction.DateCreated))
		pendingChanges := watcher.pendingChangesByTimestamp[ts]
		for _, mutatedKey := range transaction.MutatedKeys {
			mks := serializeKey(mutatedKey)
			mutatedEntity := pendingChanges[mks]
			watcher.currentEntities[mks] = mutatedEntity.Doc
			convertedEntity, err := convertSnapshotToMetaEntity(
				watcher.schema.Kinds[mutatedKey.Path[len(mutatedKey.Path)-1].Kind],
				mutatedEntity.Doc,
			)
			if err != nil {
				log.Printf("error during batch construction: %v", err)
			} else {
				batch.MutatedEntities = append(
					batch.MutatedEntities,
					convertedEntity,
				)
			}
			delete(pendingChanges, mks)
		}
		if len(pendingChanges) == 0 {
			// we can free up the map as well
			delete(watcher.pendingChangesByTimestamp, ts)
		}
	}
	for _, deletedKey := range transaction.DeletedKeys {
		dks := serializeKey(deletedKey)
		delete(watcher.currentEntities, dks)
	}

	// push the batch out
	watcher.outboundChanges <- batch

	// the transaction has been applied and can be removed from the transaction list
	return true
}

func createTransactionWatcher(ctx context.Context, client *firestore.Client, schema *Schema) (*transactionWatcher, error) {
	watcher := &transactionWatcher{
		client:                    client,
		schema:                    schema,
		currentEntities:           make(map[string]*firestore.DocumentSnapshot),
		pendingChangesByTimestamp: make(map[string]map[string]*firestore.DocumentChange),
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

			// we have transactions to process
			for len(watcher.transactions) > 0 {
				result := watcher.processTransaction(0, watcher.transactions[0])
				if result {
					// remove transaction because it's been applied
					watcher.transactions = watcher.transactions[1:]
				} else {
					// can't apply the first transaction yet, wait a little longer
					// for new entities, notifications, etc. to come in before
					// retrying
					break
				}
			}

			// release all locks
			watcher.transactionsLock.Unlock()
			watcher.pendingChangesByTimestampLock.Unlock()
			watcher.currentEntitiesLock.Unlock()
		}
	}()

	// listen for inbound changes
	go func() {
		for elem := range watcher.inboundChanges {
			watcher.pendingChangesByTimestampLock.Lock()

			ts := serializeTime(elem.Doc.UpdateTime)

			if watcher.pendingChangesByTimestamp[ts] == nil {
				watcher.pendingChangesByTimestamp[ts] = make(map[string]*firestore.DocumentChange)
			}
			watcher.pendingChangesByTimestamp[ts][serializeRef(elem.Doc.Ref)] = &elem

			watcher.pendingChangesByTimestampLock.Unlock()
		}
	}()

	// propagate outbound changes
	go func() {
		for elem := range watcher.outboundChanges {
			watcher.outboundChannelsLock.Lock()
			for _, ch := range watcher.outboundChannels {
				ch <- elem
			}
			watcher.outboundChannelsLock.Unlock()
		}
	}()

	// listen for new transactions coming in
	startTime := time.Now().Add(time.Minute * -5)
	go func() {
		transactions := watcher.client.Collection("Transaction").Where("dateSubmitted", ">=", startTime).OrderBy("dateSubmitted", firestore.Asc).Snapshots(ctx)
		for true {
			transactionSnapshot, err := transactions.Next()
			if err != nil {
				log.Printf("error during transaction watch: %v", err)
				return
			}

			watcher.transactionsLock.Lock()

			for _, change := range transactionSnapshot.Changes {
				data := change.Doc.Data()

				var mutatedRefs []*firestore.DocumentRef
				var mutatedKeys []*Key
				var deletedRefs []*firestore.DocumentRef
				var deletedKeys []*Key
				var dateSubmitted time.Time
				var description string
				if w, ok := data["mutatedKeys"].([]interface{}); ok {
					for _, ww := range w {
						if www, ok := ww.(*firestore.DocumentRef); ok {
							mutatedRefs = append(mutatedRefs, www)
						}
					}
				}
				if w, ok := data["deletedKeys"].([]interface{}); ok {
					for _, ww := range w {
						if www, ok := ww.(*firestore.DocumentRef); ok {
							deletedRefs = append(deletedRefs, www)
						}
					}
				}
				if w, ok := data["dateSubmitted"].(time.Time); ok {
					dateSubmitted = w
				}
				if w, ok := data["description"].(string); ok {
					description = w
				}
				for _, ref := range mutatedRefs {
					key, err := convertDocumentRefToMetaKey(ref)
					if err != nil {
						continue
					}
					mutatedKeys = append(mutatedKeys, key)
				}
				for _, ref := range deletedRefs {
					key, err := convertDocumentRefToMetaKey(ref)
					if err != nil {
						continue
					}
					deletedKeys = append(deletedKeys, key)
				}

				if len(mutatedKeys) == 0 && len(deletedKeys) == 0 {
					// we decoded these incorrectly, because we never write transactions
					// that don't have at least one of these
					panic("mutatedKeys and deletedKeys are both length 0!")
				}

				watcher.transactions = append(
					watcher.transactions,
					&MetaTransactionRecord{
						MutatedKeys:   mutatedKeys,
						DeletedKeys:   deletedKeys,
						DateSubmitted: convertTimeToTimestamp(dateSubmitted),
						DateCreated:   convertTimeToTimestamp(change.Doc.CreateTime),
						Description:   description,
						Id:            change.Doc.Ref.ID,
					},
				)
			}

			sort.Slice(watcher.transactions[:], func(i, j int) bool {
				if watcher.transactions[i].DateCreated.Seconds < watcher.transactions[j].DateCreated.Seconds {
					return true
				}
				if watcher.transactions[i].DateCreated.Nanos < watcher.transactions[j].DateCreated.Nanos {
					return true
				}
				return false
			})

			watcher.transactionsLock.Unlock()
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
