package main

import (
	"context"
	"fmt"
)

func (s *operationProcessor) operationUpdate(ctx context.Context, schema *Schema, req *MetaUpdateEntityRequest) (*MetaUpdateEntityResponse, error) {
	pathElements := req.Entity.Key.Path
	lastKind := pathElements[len(pathElements)-1].Kind

	var kindInfo *SchemaKind
	for kindName, kind := range schema.Kinds {
		if kindName == lastKind {
			kindInfo = kind
			break
		}
	}
	if kindInfo == nil {
		return nil, fmt.Errorf("no such kind")
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

	_, err = ref.Set(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("can't set data against entity (Firestore): %v", err)
	}

	return &MetaUpdateEntityResponse{
		Entity: req.Entity,
	}, nil
}
