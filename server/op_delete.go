package main

import (
	"context"
	"fmt"
)

func (s *operationProcessor) operationDelete(ctx context.Context, schema *Schema, req *MetaDeleteEntityRequest) (*MetaDeleteEntityResponse, error) {
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

	_, err = ref.Delete(ctx)
	if err != nil {
		return nil, err
	}

	response := &MetaDeleteEntityResponse{
		Entity: entity,
	}

	return response, nil
}
