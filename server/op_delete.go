package main

import (
	"context"

	"cloud.google.com/go/firestore"
)

func (s *operationProcessor) operationDeleteRead(ctx context.Context, schema *Schema, req *MetaDeleteEntityRequest) (interface{}, error) {
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

	return snapshot, nil
}

func (s *operationProcessor) operationDeleteWrite(ctx context.Context, schema *Schema, req *MetaDeleteEntityRequest, readState interface{}) (*MetaDeleteEntityResponse, error) {
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

	snapshot := readState.(*firestore.DocumentSnapshot)

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
