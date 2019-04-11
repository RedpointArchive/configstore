package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

func operationList(ctx context.Context, schema *Schema, req *MetaListEntitiesRequest) (*MetaListEntitiesResponse, error) {
	var start interface{}
	if req.Start != nil {
		if len(req.Start[:]) > 0 {
			start = string(req.Start[:])
		}
	}

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

	var err error
	var snapshots []*firestore.DocumentSnapshot
	if (req.Limit == 0) && start == nil {
		snapshots, err = client.Collection(req.KindName).Documents(ctx).GetAll()
	} else if req.Limit == 0 {
		snapshots, err = client.Collection(req.KindName).OrderBy(firestore.DocumentID, firestore.Asc).StartAfter(start.(string)).Documents(ctx).GetAll()
	} else if start == nil {
		snapshots, err = client.Collection(req.KindName).Limit(int(req.Limit)).Documents(ctx).GetAll()
	} else {
		snapshots, err = client.Collection(req.KindName).OrderBy(firestore.DocumentID, firestore.Asc).StartAfter(start.(string)).Limit(int(req.Limit)).Documents(ctx).GetAll()
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
