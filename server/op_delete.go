package main

import (
	"context"
)

func (s *operationProcessor) operationDelete(ctx context.Context, schema *Schema, req *MetaDeleteEntityRequest) (*MetaDeleteEntityResponse, error) {
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

	err = s.tx.Delete(ref)
	if err != nil {
		return nil, err
	}

	response := &MetaDeleteEntityResponse{
		Entity: entity,
	}

	return response, nil
}
