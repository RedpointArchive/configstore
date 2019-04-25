package main

import (
	"context"
	"fmt"
)

func (s *operationProcessor) operationUpdateRead(ctx context.Context, schema *Schema, req *MetaUpdateEntityRequest) (interface{}, error) {
	return nil, nil
}

func (s *operationProcessor) operationUpdateWrite(ctx context.Context, schema *Schema, req *MetaUpdateEntityRequest, readState interface{}) (*MetaUpdateEntityResponse, error) {
	pathElements := req.Entity.Key.Path
	lastKind := pathElements[len(pathElements)-1].Kind

	kindInfo, err := findSchemaKindByName(schema, lastKind)
	if err != nil {
		return nil, err
	}

	ref, data, err := convertMetaEntityToRefAndDataMap(
		s.client,
		req.Entity,
		kindInfo,
	)
	if err != nil {
		return nil, fmt.Errorf("can't convert meta entity to ref and map: %v", err)
	}

	if ref == nil {
		return nil, fmt.Errorf("entity must be set")
	}

	err = s.tx.Set(ref, data)
	if err != nil {
		return nil, fmt.Errorf("can't set data against entity (Firestore): %v", err)
	}

	return &MetaUpdateEntityResponse{
		Entity: req.Entity,
	}, nil
}
