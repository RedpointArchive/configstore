package main

import (
	"context"
)

func (s *operationProcessor) operationGetRead(ctx context.Context, schema *Schema, req *MetaGetEntityRequest) (interface{}, error) {
	kindInfo, err := findSchemaKindByName(schema, req.KindName)
	if err != nil {
		return nil, err
	}

	ref, err := convertMetaKeyToDocumentRef(
		s.client,
		req.Key,
	)
	if err != nil {
		return nil, err
	}

	snapshot, err := s.tx.Get(ref)
	if err != nil {
		return nil, err
	}

	entity, err := convertSnapshotToMetaEntity(kindInfo, snapshot)
	if err != nil {
		return nil, err
	}

	response := &MetaGetEntityResponse{
		Entity: entity,
	}

	return response, nil
}

func (s *operationProcessor) operationGetWrite(ctx context.Context, schema *Schema, req *MetaGetEntityRequest, readState interface{}) (*MetaGetEntityResponse, error) {
	return readState.(*MetaGetEntityResponse), nil
}
