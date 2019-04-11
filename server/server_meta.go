package main

import (
	"context"
)

type configstoreMetaServiceServer struct {
	schema *Schema
}

func (s *configstoreMetaServiceServer) GetSchema(ctx context.Context, req *GetSchemaRequest) (*GetSchemaResponse, error) {
	return &GetSchemaResponse{
		Schema: s.schema,
	}, nil
}

func (s *configstoreMetaServiceServer) GetDefaultPartitionId(ctx context.Context, req *GetDefaultPartitionIdRequest) (*GetDefaultPartitionIdResponse, error) {
	firestoreTestCollection := client.Collection("Test")
	firestoreNamespace := firestoreTestCollection.Path[0:(len(firestoreTestCollection.Path) - len(firestoreTestCollection.ID) - 1)]

	return &GetDefaultPartitionIdResponse{
		Namespace: firestoreNamespace,
	}, nil
}

func (s *configstoreMetaServiceServer) MetaList(ctx context.Context, req *MetaListEntitiesRequest) (*MetaListEntitiesResponse, error) {
	resp, err := processTransaction(
		ctx,
		s.schema,
		&MetaTransaction{
			Operations: []*MetaOperation{
				&MetaOperation{
					Operation: req,
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
	resp, err := processTransaction(
		ctx,
		s.schema,
		&MetaTransaction{
			Operations: []*MetaOperation{
				&MetaOperation{
					Operation: req,
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
	resp, err := processTransaction(
		ctx,
		s.schema,
		&MetaTransaction{
			Operations: []*MetaOperation{
				&MetaOperation{
					Operation: req,
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
	resp, err := processTransaction(
		ctx,
		s.schema,
		&MetaTransaction{
			Operations: []*MetaOperation{
				&MetaOperation{
					Operation: req,
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
	resp, err := processTransaction(
		ctx,
		s.schema,
		&MetaTransaction{
			Operations: []*MetaOperation{
				&MetaOperation{
					Operation: req,
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
