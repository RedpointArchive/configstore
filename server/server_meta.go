package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
)

type configstoreMetaServiceServer struct {
	firestoreClient      *firestore.Client
	schema               *Schema
	transactionProcessor *transactionProcessor
	transactionWatcher   *transactionWatcher
}

func createConfigstoreMetaServiceServer(
	firestoreClient *firestore.Client,
	schema *Schema,
	transactionProcessor *transactionProcessor,
	transactionWatcher *transactionWatcher,
) *configstoreMetaServiceServer {
	return &configstoreMetaServiceServer{
		firestoreClient:      firestoreClient,
		schema:               schema,
		transactionProcessor: transactionProcessor,
		transactionWatcher:   transactionWatcher,
	}
}

func (s *configstoreMetaServiceServer) GetSchema(ctx context.Context, req *GetSchemaRequest) (*GetSchemaResponse, error) {
	return &GetSchemaResponse{
		Schema: s.schema,
	}, nil
}

func (s *configstoreMetaServiceServer) GetDefaultPartitionId(ctx context.Context, req *GetDefaultPartitionIdRequest) (*GetDefaultPartitionIdResponse, error) {
	firestoreTestCollection := s.firestoreClient.Collection("Test")
	firestoreNamespace := firestoreTestCollection.Path[0:(len(firestoreTestCollection.Path) - len(firestoreTestCollection.ID) - 1)]

	return &GetDefaultPartitionIdResponse{
		Namespace: firestoreNamespace,
	}, nil
}

func (s *configstoreMetaServiceServer) MetaList(ctx context.Context, req *MetaListEntitiesRequest) (*MetaListEntitiesResponse, error) {
	resp, err := s.transactionProcessor.processTransaction(
		ctx,
		s.schema,
		&MetaTransaction{
			Operations: []*MetaOperation{
				&MetaOperation{
					Operation: &MetaOperation_ListRequest{
						ListRequest: req,
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.OperationResults[0].Error != nil {
		return nil, fmt.Errorf("%s", resp.OperationResults[0].Error.ErrorMessage)
	}
	return resp.OperationResults[0].GetListResponse(), nil
}

func (s *configstoreMetaServiceServer) MetaGet(ctx context.Context, req *MetaGetEntityRequest) (*MetaGetEntityResponse, error) {
	resp, err := s.transactionProcessor.processTransaction(
		ctx,
		s.schema,
		&MetaTransaction{
			Operations: []*MetaOperation{
				&MetaOperation{
					Operation: &MetaOperation_GetRequest{
						GetRequest: req,
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.OperationResults[0].Error != nil {
		return nil, fmt.Errorf("%s", resp.OperationResults[0].Error.ErrorMessage)
	}
	return resp.OperationResults[0].GetGetResponse(), nil
}

func (s *configstoreMetaServiceServer) MetaUpdate(ctx context.Context, req *MetaUpdateEntityRequest) (*MetaUpdateEntityResponse, error) {
	resp, err := s.transactionProcessor.processTransaction(
		ctx,
		s.schema,
		&MetaTransaction{
			Operations: []*MetaOperation{
				&MetaOperation{
					Operation: &MetaOperation_UpdateRequest{
						UpdateRequest: req,
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.OperationResults[0].Error != nil {
		return nil, fmt.Errorf("%s", resp.OperationResults[0].Error.ErrorMessage)
	}
	return resp.OperationResults[0].GetUpdateResponse(), nil
}

func (s *configstoreMetaServiceServer) MetaDelete(ctx context.Context, req *MetaDeleteEntityRequest) (*MetaDeleteEntityResponse, error) {
	resp, err := s.transactionProcessor.processTransaction(
		ctx,
		s.schema,
		&MetaTransaction{
			Operations: []*MetaOperation{
				&MetaOperation{
					Operation: &MetaOperation_DeleteRequest{
						DeleteRequest: req,
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.OperationResults[0].Error != nil {
		return nil, fmt.Errorf("%s", resp.OperationResults[0].Error.ErrorMessage)
	}
	return resp.OperationResults[0].GetDeleteResponse(), nil
}

func (s *configstoreMetaServiceServer) MetaCreate(ctx context.Context, req *MetaCreateEntityRequest) (*MetaCreateEntityResponse, error) {
	resp, err := s.transactionProcessor.processTransaction(
		ctx,
		s.schema,
		&MetaTransaction{
			Operations: []*MetaOperation{
				&MetaOperation{
					Operation: &MetaOperation_CreateRequest{
						CreateRequest: req,
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.OperationResults[0].Error != nil {
		return nil, fmt.Errorf("%s", resp.OperationResults[0].Error.ErrorMessage)
	}
	return resp.OperationResults[0].GetCreateResponse(), nil
}

func (s *configstoreMetaServiceServer) ApplyTransaction(ctx context.Context, req *MetaTransaction) (*MetaTransactionResult, error) {
	resp, err := s.transactionProcessor.processTransaction(
		ctx,
		s.schema,
		req,
	)
	return resp, err
}

func (s *configstoreMetaServiceServer) WatchTransactions(req *WatchTransactionsRequest, srv ConfigstoreMetaService_WatchTransactionsServer) error {
	ch := make(chan *MetaTransactionBatch)
	s.transactionWatcher.RegisterChannel(ch)
	defer s.transactionWatcher.DeregisterChannel(ch)

	connected := true
	for connected {
		select {
		case msg := <-ch:
			srv.Send(&WatchTransactionsResponse{
				Batch: msg,
			})
		case <-time.After(1 * time.Second):
			err := srv.Context().Err()
			if err != nil {
				fmt.Println("detected client disconnection")
				connected = false
				break
			}
		}
	}

	return nil
}
