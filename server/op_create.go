package main

import (
	"context"
)

func (s *operationProcessor) operationCreateRead(ctx context.Context, schema *Schema, req *MetaCreateEntityRequest) (interface{}, error) {
	return nil, nil
}

func (s *operationProcessor) operationCreateWrite(ctx context.Context, schema *Schema, req *MetaCreateEntityRequest, readState interface{}) (*MetaCreateEntityResponse, error) {
	kindInfo, err := findSchemaKindByName(schema, req.KindName)
	if err != nil {
		return nil, err
	}

	if req.Entity.Key == nil {
		// we need to automatically generate a key for this entity
		firestoreCollection := s.client.Collection(req.KindName)

		newKey, err := convertDocumentRefToMetaKey(
			firestoreCollection.NewDoc(),
		)
		if err != nil {
			return nil, err
		}

		req.Entity.Key = newKey
	}

	ref, data, err := convertMetaEntityToRefAndDataMap(
		s.client,
		req.Entity,
		kindInfo,
	)
	if err != nil {
		return nil, err
	}

	if ref.ID == "" {
		ref = ref.Parent.NewDoc()
	}
	err = s.tx.Create(ref, data)
	if err != nil {
		return nil, err
	}

	key, err := convertDocumentRefToMetaKey(
		ref,
	)
	if err != nil {
		return nil, err
	}

	// set the key
	req.Entity.Key = key

	return &MetaCreateEntityResponse{
		Entity: req.Entity,
	}, nil
}
