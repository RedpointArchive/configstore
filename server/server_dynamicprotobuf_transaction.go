package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jhump/protoreflect/desc/builder"
	"github.com/jhump/protoreflect/dynamic"

	"google.golang.org/grpc"

	"cloud.google.com/go/firestore"
)

type configstoreDynamicProtobufTransactionService struct {
	firestoreClient      *firestore.Client
	genResult            *generatorResult
	service              *builder.ServiceBuilder
	schema               *Schema
	transactionProcessor *transactionProcessor
	transactionWatcher   *transactionWatcher
}

func createConfigstoreDynamicProtobufTransactionServer(
	firestoreClient *firestore.Client,
	genResult *generatorResult,
	service *builder.ServiceBuilder,
	schema *Schema,
	transactionWatcher *transactionWatcher,
) *configstoreDynamicProtobufTransactionService {
	return &configstoreDynamicProtobufTransactionService{
		firestoreClient:      firestoreClient,
		genResult:            genResult,
		service:              service,
		schema:               schema,
		transactionProcessor: createTransactionProcessor(firestoreClient),
		transactionWatcher:   transactionWatcher,
	}
}

func (s *configstoreDynamicProtobufTransactionService) getMetaServiceServer() *configstoreMetaServiceServer {
	return createConfigstoreMetaServiceServer(
		s.firestoreClient,
		s.schema,
		s.transactionProcessor,
		s.transactionWatcher,
	)
}

func (s *configstoreDynamicProtobufTransactionService) dynamicProtobufTransactionWatch(ctx context.Context, srv interface{}, stream grpc.ServerStream) error {
	messageFactory := dynamic.NewMessageFactoryWithDefaults()

	if !s.transactionWatcher.isConsistent {
		return fmt.Errorf("configstore is not yet transactionally consistent because it is starting up, please try again in a moment")
	}

	// lock before registering for notifications, so we don't miss any transactions
	// that get applied to the entity state
	s.transactionWatcher.CurrentEntitiesTakeReadLock()
	hasReadLock := true
	releaseLockIfHeld := func() {
		if hasReadLock {
			s.transactionWatcher.CurrentEntitiesReleaseReadLock()
			hasReadLock = false
		}
	}
	defer releaseLockIfHeld()

	// register for new transaction notifications
	ch := make(chan *MetaTransactionBatch)
	s.transactionWatcher.RegisterChannel(ch)
	defer s.transactionWatcher.DeregisterChannel(ch)

	// send down the initial state of the database
	initialState := messageFactory.NewDynamicMessage(s.genResult.MessageMap["TypedTransactionInitialState"])
	var transactionEntities []*dynamic.Message
	for _, snapshot := range s.transactionWatcher.currentEntities {
		key, err := convertDocumentRefToMetaKey(snapshot.Ref)
		if err != nil {
			return err
		}
		kindName := key.Path[len(key.Path)-1].Kind
		kind := s.transactionWatcher.schema.Kinds[kindName]
		metaEntity, err := convertSnapshotToMetaEntity(
			kind,
			snapshot,
		)
		if err != nil {
			return err
		}
		entityMessage, err := convertMetaEntityToDynamicMessage(
			messageFactory,
			s.genResult.MessageMap[kindName],
			metaEntity,
			s.genResult.CommonMessageDescriptors,
			kind,
		)
		if err != nil {
			return err
		}
		transactionEntity := messageFactory.NewDynamicMessage(s.genResult.MessageMap["TypedTransactionEntity"])
		transactionEntity.SetFieldByNumber(
			int(kind.Id),
			entityMessage,
		)
		transactionEntities = append(
			transactionEntities,
			transactionEntity,
		)
	}
	initialState.SetFieldByName("entities", transactionEntities)
	releaseLockIfHeld()

	watchTransactionResponse := messageFactory.NewDynamicMessage(s.genResult.MessageMap["TypedWatchTransactionsResponse"])
	watchTransactionResponse.SetFieldByName("initialState", initialState)
	stream.SendMsg(watchTransactionResponse)

	// send down transactions as they occur
	connected := true
	for connected {
		select {
		case msg := <-ch:
			var mutatedEntities []*dynamic.Message
			for _, mutatedEntity := range msg.MutatedEntities {
				kindName := mutatedEntity.Key.Path[len(mutatedEntity.Key.Path)-1].Kind
				kind := s.transactionWatcher.schema.Kinds[kindName]
				typedMutatedEntity, err := convertMetaEntityToDynamicMessage(
					messageFactory,
					s.genResult.MessageMap[kindName],
					mutatedEntity,
					s.genResult.CommonMessageDescriptors,
					kind,
				)
				if err != nil {
					return err
				}
				transactionEntity := messageFactory.NewDynamicMessage(s.genResult.MessageMap["TypedTransactionEntity"])
				transactionEntity.SetFieldByNumber(
					int(kind.Id),
					typedMutatedEntity,
				)
				mutatedEntities = append(
					mutatedEntities,
					transactionEntity,
				)
			}

			batch := messageFactory.NewDynamicMessage(s.genResult.MessageMap["TypedTransactionBatch"])
			batch.SetFieldByName("mutatedEntities", mutatedEntities)
			batch.SetFieldByName("deletedKeys", msg.DeletedKeys)
			batch.SetFieldByName("description", msg.Description)
			batch.SetFieldByName("id", msg.Id)

			watchTransactionResponse := messageFactory.NewDynamicMessage(s.genResult.MessageMap["TypedWatchTransactionsResponse"])
			watchTransactionResponse.SetFieldByName("batch", batch)
			stream.SendMsg(watchTransactionResponse)
		case <-time.After(1 * time.Second):
			err := stream.Context().Err()
			if err != nil {
				connected = false
				break
			}
		}
	}

	return nil
}
