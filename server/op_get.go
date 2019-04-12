package main

import (
	"context"
)

func (s *operationProcessor) operationGet(ctx context.Context, schema *Schema, req *MetaGetEntityRequest) (*MetaGetEntityResponse, error) {
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

	snapshot, err := ref.Get(ctx)
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
