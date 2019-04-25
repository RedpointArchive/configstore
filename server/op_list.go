package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

func (s *operationProcessor) operationListRead(ctx context.Context, schema *Schema, req *MetaListEntitiesRequest) (interface{}, error) {
	var start interface{}
	if req.Start != nil {
		if len(req.Start[:]) > 0 {
			start = string(req.Start[:])
		}
	}

	kindInfo, err := findSchemaKindByName(schema, req.KindName)
	if err != nil {
		return nil, err
	}

	var snapshots []*firestore.DocumentSnapshot
	if runWithoutFirestoreTransactionalQueries() {
		if (req.Limit == 0) && start == nil {
			snapshots, err = s.client.Collection(req.KindName).Documents(ctx).GetAll()
		} else if req.Limit == 0 {
			snapshots, err = s.client.Collection(req.KindName).OrderBy(firestore.DocumentID, firestore.Asc).StartAfter(start.(string)).Documents(ctx).GetAll()
		} else if start == nil {
			snapshots, err = s.client.Collection(req.KindName).Limit(int(req.Limit)).Documents(ctx).GetAll()
		} else {
			snapshots, err = s.client.Collection(req.KindName).OrderBy(firestore.DocumentID, firestore.Asc).StartAfter(start.(string)).Limit(int(req.Limit)).Documents(ctx).GetAll()
		}
	} else {
		if (req.Limit == 0) && start == nil {
			snapshots, err = s.tx.Documents(s.client.Collection(req.KindName)).GetAll()
		} else if req.Limit == 0 {
			snapshots, err = s.tx.Documents(s.client.Collection(req.KindName).OrderBy(firestore.DocumentID, firestore.Asc).StartAfter(start.(string))).GetAll()
		} else if start == nil {
			snapshots, err = s.tx.Documents(s.client.Collection(req.KindName).Limit(int(req.Limit))).GetAll()
		} else {
			snapshots, err = s.tx.Documents(s.client.Collection(req.KindName).OrderBy(firestore.DocumentID, firestore.Asc).StartAfter(start.(string)).Limit(int(req.Limit))).GetAll()
		}
	}

	if err != nil {
		return nil, err
	}

	var entities []*MetaEntity
	for _, snapshot := range snapshots {
		entity, err := convertSnapshotToMetaEntity(kindInfo, snapshot)
		if err != nil {
			fmt.Printf("%v", err)
			continue
		}
		entities = append(entities, entity)
	}

	response := &MetaListEntitiesResponse{
		Entities: entities,
	}

	if !(req.Limit == 0) {
		if uint32(len(entities)) < req.Limit {
			response.MoreResults = false
		} else {
			// TODO: query to see if there really are more results, to make this behave like datastore
			response.MoreResults = true
			last := snapshots[len(snapshots)-1]
			response.Next = []byte(last.Ref.ID)
		}
	} else {
		response.MoreResults = false
	}

	return response, nil
}

func (s *operationProcessor) operationListWrite(ctx context.Context, schema *Schema, req *MetaListEntitiesRequest, readState interface{}) (*MetaListEntitiesResponse, error) {
	return readState.(*MetaListEntitiesResponse), nil
}
