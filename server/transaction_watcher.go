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

	initialReadTimeByKind map[string]time.Time

	isConsistent bool

	waitTransactionCount int
}

func (watcher *transactionWatcher) CurrentEntitiesTakeReadLock() {
	watcher.currentEntitiesLock.RLock()
}
func (watcher *transactionWatcher) CurrentEntitiesReleaseReadLock() {
	watcher.currentEntitiesLock.RUnlock()
}
func (watcher *transactionWatcher) CurrentEntitiesTakeWriteLock() {
	watcher.currentEntitiesLock.Lock()
}
func (watcher *transactionWatcher) CurrentEntitiesReleaseWriteLock() {
	watcher.currentEntitiesLock.Unlock()
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
	close(oldCh)

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
	// at this point, we know the transaction hasn't already been applied to
	// current entities, so we need to check the pendingChangesByTimestamp
	// to see if we have all the mutated keys present (deleted keys won't exist;
	// we just remove them from current entities once we have all the mutations
	// we need)
	if len(transaction.MutatedKeys) > 0 {
		ts := serializeTime(convertTimestampToTime(transaction.DateCreated))
		pendingChanges, ok := watcher.pendingChangesByTimestamp[ts]
		if !ok {
			// auto-create map for simplicity
			watcher.pendingChangesByTimestamp[ts] = make(map[string]*firestore.DocumentChange)
			pendingChanges = watcher.pendingChangesByTimestamp[ts]
		}
		for _, mutatedKey := range transaction.MutatedKeys {
			mks := serializeKey(mutatedKey)

			// we have a mutated key, check if the initial read for this entity kind occurred after
			// this transaction was created - if it did, we may never see the pending arrive because
			// the entity has since been deleted (and we already have the up to date state of this type
			// of change anyway).
			kindName := mutatedKey.Path[len(mutatedKey.Path)-1].Kind
			if watcher.initialReadTimeByKind[kindName].After(convertTimestampToTime(transaction.DateCreated)) {
				// we won't receive a pending change for this entity, pretend that the "current entity" state is the
				// pending change we're waiting for to make further logic work later
				if doc, ok := watcher.currentEntities[mks]; ok {
					watcher.pendingChangesByTimestamp[ts][mks] = &firestore.DocumentChange{
						Kind:     firestore.DocumentAdded,
						Doc:      doc,
						OldIndex: -1,
						NewIndex: -1,
					}
				} else {
					watcher.pendingChangesByTimestamp[ts][mks] = &firestore.DocumentChange{
						Kind:     firestore.DocumentRemoved,
						Doc:      nil,
						OldIndex: -1,
						NewIndex: -1,
					}
				}
			} else {
				// our initial state for entities of this kind didn't include this entity (and therefore we can't
				// have had a newer "delete" transaction prevent the entity from appearing, so we need to check it)
				_, ok := pendingChanges[mks]
				if !ok {
					// we don't have this entity's snapshot yet
					fmt.Printf("%s: can't process transaction, waiting on entity snapshot with key: %s\n", transaction.Id, mks)
					watcher.waitTransactionCount++
					if watcher.waitTransactionCount > 30 {
						// we aren't making progress - we've stalled for some reason (race condition, etc.)
						// this indicates a bug in configstore, but we need to make sure applications using configstore
						// can recover while the issue gets resolved, so what we do is panic() here and expect the
						// orchestrator (such as Kubernetes) to restart the pod. since you should also be deploying
						// configstore with more than 1 replica, this should result in no downtime for applications
						panic(fmt.Sprintf("panic! %s: unable to make progress on transaction, still waiting on entity snapshot with key: %s\n", transaction.Id, mks))
					}
					return false
				}
			}
		}
	}

	// at this point, we have all the changes for this transaction in pendingChanges,
	// we can apply them to currentEntities
	batch := &MetaTransactionBatch{
		Id:              transaction.Id,
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
			// doc can be nil here if we're recovering from a partial read during
			// initial load, and we can't provide the historical version of a snapshot
			// after it's been deleted
			//
			// clients can't connect to configstore's WatchTransactions endpoint
			// until the watcher is "consistent", which is only true once 30 seconds have
			// passed since reading all of the initial entities (this gives enough time
			// for any transactions occurring during startup to come in) and once all pending
			// transactions have been processed. thus we can reasonably be certain that
			// configstore's in memory version of the database is consistent, and all future
			// transactions will be applied atomically and consistently via this code (since once
			// configstore has started, we *can* see historical versions of snapshots that are deleted).
			if mutatedEntity.Doc != nil {
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
	watcher.waitTransactionCount = 0

	// the transaction has been applied and can be removed from the transaction list
	fmt.Printf("%s: finished processing transaction (%d transactions left to process)\n", transaction.Id, len(watcher.transactions)-1)
	return true
}

func safeSendBatchToChannel(ch chan *MetaTransactionBatch, batch *MetaTransactionBatch) bool {
	defer func() {
		recover()
	}()
	ch <- batch
	return true
}

func createTransactionWatcher(ctx context.Context, client *firestore.Client, schema *Schema) (*transactionWatcher, error) {
	watcher := &transactionWatcher{
		client:                    client,
		schema:                    schema,
		currentEntities:           make(map[string]*firestore.DocumentSnapshot),
		pendingChangesByTimestamp: make(map[string]map[string]*firestore.DocumentChange),
		inboundChanges:            make(chan firestore.DocumentChange),
		outboundChanges:           make(chan *MetaTransactionBatch),
		transactions:              nil,
		initialReadTimeByKind:     make(map[string]time.Time),
		isConsistent:              false,
	}

	watcher.CurrentEntitiesTakeWriteLock()
	defer watcher.CurrentEntitiesReleaseWriteLock()

	// process transactions every second
	go func() {
		// prevents transaction processing from starting until
		// after we have got the initial reads of all entities (so
		// watcher.initialReadTimeByKind is populated)
		watcher.CurrentEntitiesTakeWriteLock()
		watcher.CurrentEntitiesReleaseWriteLock()

		for true {
			// wait a second
			time.Sleep(time.Second * 1)

			// preemptive lock to check if we have any transactions to process at all
			watcher.transactionsLock.RLock()
			if len(watcher.transactions) == 0 {
				if !watcher.isConsistent {
					deadline := time.Now().Add(time.Second * 30)
					isConsistent := true
					for _, readTime := range watcher.initialReadTimeByKind {
						if !readTime.Before(deadline) {
							isConsistent = false
							break
						}
					}
					if isConsistent {
						fmt.Printf("configstore is now consistent and ready to serve transactions\n")
						watcher.isConsistent = isConsistent
					}
				}
				watcher.transactionsLock.RUnlock()
				continue
			}
			watcher.transactionsLock.RUnlock()

			// obtain all locks
			watcher.transactionsLock.Lock()
			watcher.pendingChangesByTimestampLock.Lock()
			watcher.CurrentEntitiesTakeWriteLock()

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
					fmt.Printf("%s: unable to process yet\n", watcher.transactions[0].Id)
					break
				}
			}

			// release all locks
			watcher.CurrentEntitiesReleaseWriteLock()
			watcher.pendingChangesByTimestampLock.Unlock()
			watcher.transactionsLock.Unlock()
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
				safeSendBatchToChannel(ch, elem)
			}
			watcher.outboundChannelsLock.Unlock()
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

	if runWithoutFirestoreTransactionalQueries() {
		// for each kind, fill in the current entities
		for kindName := range schema.Kinds {
			documents, err := watcher.client.Collection(kindName).Documents(ctx).GetAll()
			if err != nil {
				return nil, err
			}
			for _, document := range documents {
				watcher.initialReadTimeByKind[kindName] = document.ReadTime
				watcher.currentEntities[serializeRef(document.Ref)] = document
			}
		}
	} else {
		err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			// for each kind, fill in the current entities
			for kindName := range schema.Kinds {
				documents, err := tx.Documents(watcher.client.Collection(kindName)).GetAll()
				if err != nil {
					return err
				}
				for _, document := range documents {
					watcher.initialReadTimeByKind[kindName] = document.ReadTime
					watcher.currentEntities[serializeRef(document.Ref)] = document
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	// listen for new transactions coming in
	go func() {
		transactions := watcher.client.Collection("Transaction").Where("dateSubmitted", ">=", time.Now().Add(time.Second*-60)).OrderBy("dateSubmitted", firestore.Asc).Snapshots(ctx)
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

				fmt.Printf("%s: transaction arrived\n", change.Doc.Ref.ID)
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

			sort.Slice(watcher.transactions, func(i, j int) bool {
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

	return watcher, nil
}
