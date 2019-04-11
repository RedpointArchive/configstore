package main

import (
	"context"
	"fmt"
)

func (s *operationProcessor) operationGet(ctx context.Context, schema *Schema, req *MetaGetEntityRequest) (*MetaGetEntityResponse, error) {
	var kindInfo *SchemaKind
	for kindName, kind := range schema.Kinds {
		if kindName == req.KindName {
			kindInfo = kind
			break
		}
	}
	if kindInfo == nil {
		return nil, fmt.Errorf("no such kind")
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
