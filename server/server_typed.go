package main

import (
	"context"
)

type configstoreTypedService struct {
	GenResult *generatorResult
	Service   *builder.ServiceBuilder
	KindName  string
}

func (s *configstoreTypedService) TypedList(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	messageFactory := dynamic.NewMessageFactoryWithDefaults()

	requestMessageDescriptor := s.GenResult.MessageMap[fmt.Sprintf("List%sRequest", s.KindName)]
	in := messageFactory.NewDynamicMessage(requestMessageDescriptor)
	if err := dec(in); err != nil {
		return nil, err
	}

	startBytes, err := in.TryGetFieldByName("start")
	if err != nil {
		return nil, err
	}
	limit, err := in.TryGetFieldByName("limit")
	if err != nil {
		return nil, err
	}

	var start interface{}
	if startBytes != nil {
		if len(startBytes.([]byte)[:]) > 0 {
			start = string(startBytes.([]byte)[:])
		}
	}

	var snapshots []*firestore.DocumentSnapshot
	if (limit == nil || limit.(uint32) == 0) && start == nil {
		snapshots, err = client.Collection(s.KindName).Documents(ctx).GetAll()
	} else if limit == nil || limit.(uint32) == 0 {
		snapshots, err = client.Collection(s.KindName).OrderBy(firestore.DocumentID, firestore.Asc).StartAfter(start.(string)).Documents(ctx).GetAll()
	} else if start == nil {
		snapshots, err = client.Collection(s.KindName).Limit(int(limit.(uint32))).Documents(ctx).GetAll()
	} else {
		snapshots, err = client.Collection(s.KindName).OrderBy(firestore.DocumentID, firestore.Asc).StartAfter(start.(string)).Limit(int(limit.(uint32))).Documents(ctx).GetAll()
	}

	if err != nil {
		return nil, err
	}

	var entities []*dynamic.Message
	for _, snapshot := range snapshots {
		entity, err := convertSnapshotToDynamicMessage(
			messageFactory,
			s.GenResult.MessageMap[s.KindName],
			snapshot,
			s.GenResult.CommonMessageDescriptors,
		)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	responseMessageDescriptor := s.GenResult.MessageMap[fmt.Sprintf("List%sResponse", s.KindName)]
	out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
	out.SetFieldByName("entities", entities)

	if !(limit == nil || limit.(uint32) == 0) {
		if uint32(len(entities)) < limit.(uint32) {
			out.SetFieldByName("moreResults", false)
		} else {
			// TODO: query to see if there really are more results, to make this behave like datastore
			out.SetFieldByName("moreResults", true)
			last := snapshots[len(snapshots)-1]
			out.SetFieldByName("next", []byte(last.Ref.ID))
		}
	} else {
		out.SetFieldByName("moreResults", false)
	}

	return out, nil
}

func (s *configstoreTypedService) TypedGet(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	messageFactory := dynamic.NewMessageFactoryWithDefaults()

	requestMessageDescriptor := s.GenResult.MessageMap[fmt.Sprintf("Get%sRequest", s.KindName)]
	in := messageFactory.NewDynamicMessage(requestMessageDescriptor)
	if err := dec(in); err != nil {
		return nil, err
	}

	key, err := in.TryGetFieldByName("key")
	if err != nil {
		return nil, err
	}

	keyV, ok := key.(*dynamic.Message)
	if !ok {
		return nil, fmt.Errorf("unable to read key")
	}

	ref, err := convertKeyToDocumentRef(
		client,
		keyV,
	)
	if err != nil {
		return nil, err
	}

	// TODO: Validate that the last component of Kind in the DocumentRef
	// matches our expected type.

	snapshot, err := ref.Get(ctx)
	if err != nil {
		return nil, err
	}

	entity, err := convertSnapshotToDynamicMessage(
		messageFactory,
		s.GenResult.MessageMap[s.KindName],
		snapshot,
		s.GenResult.CommonMessageDescriptors,
	)
	if err != nil {
		return nil, err
	}

	responseMessageDescriptor := s.GenResult.MessageMap[fmt.Sprintf("Get%sResponse", s.KindName)]
	out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
	out.SetFieldByName("entity", entity)

	return out, nil
}

func (s *configstoreTypedService) TypedUpdate(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	messageFactory := dynamic.NewMessageFactoryWithDefaults()

	requestMessageDescriptor := s.GenResult.MessageMap[fmt.Sprintf("Update%sRequest", s.KindName)]
	in := messageFactory.NewDynamicMessage(requestMessageDescriptor)
	if err := dec(in); err != nil {
		return nil, err
	}

	entity, err := in.TryGetFieldByName("entity")
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, fmt.Errorf("entity must not be nil")
	}

	// Get the existing version so we can make sure we're not updating
	// read-only fields.
	keyRaw, err := entity.(*dynamic.Message).TryGetFieldByName("key")
	if err != nil {
		return nil, err
	}

	keyCon, ok := keyRaw.(*dynamic.Message)
	if !ok {
		return nil, fmt.Errorf("key of unexpected type")
	}

	ref, err := convertKeyToDocumentRef(
		client,
		keyCon,
	)
	if err != nil {
		return nil, err
	}

	snapshot, err := ref.Get(ctx)
	if err != nil {
		return nil, err
	}

	ref, data, err := convertDynamicMessageIntoRefAndDataMap(
		client,
		messageFactory,
		s.GenResult.MessageMap[s.KindName],
		entity.(*dynamic.Message),
		snapshot,
		s.GenResult.Schema.Kinds[s.KindName],
	)
	if err != nil {
		return nil, err
	}

	if ref == nil {
		return nil, fmt.Errorf("entity must be set")
	}

	_, err = ref.Set(ctx, data)
	if err != nil {
		return nil, err
	}

	responseMessageDescriptor := s.GenResult.MessageMap[fmt.Sprintf("Update%sResponse", s.KindName)]
	out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
	out.SetFieldByName("entity", entity)

	return out, nil
}

func (s *configstoreTypedService) TypedCreate(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	messageFactory := dynamic.NewMessageFactoryWithDefaults()

	requestMessageDescriptor := s.GenResult.MessageMap[fmt.Sprintf("Create%sRequest", s.KindName)]
	in := messageFactory.NewDynamicMessage(requestMessageDescriptor)
	if err := dec(in); err != nil {
		return nil, err
	}

	entity, err := in.TryGetFieldByName("entity")
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, fmt.Errorf("entity must not be nil")
	}

	ref, data, err := convertDynamicMessageIntoRefAndDataMap(
		client,
		messageFactory,
		s.GenResult.MessageMap[s.KindName],
		entity.(*dynamic.Message),
		nil,
		s.GenResult.Schema.Kinds[s.KindName],
	)
	if err != nil {
		return nil, err
	}

	if ref.ID == "" {
		ref, _, err = ref.Parent.Add(ctx, data)
	} else {
		_, err = ref.Create(ctx, data)
	}
	if err != nil {
		return nil, err
	}

	key, err := convertDocumentRefToKey(
		messageFactory,
		ref,
		s.GenResult.CommonMessageDescriptors,
	)
	if err != nil {
		return nil, err
	}

	// set the ID back
	entity.(*dynamic.Message).SetFieldByName("key", key)

	responseMessageDescriptor := s.GenResult.MessageMap[fmt.Sprintf("Create%sResponse", s.KindName)]
	out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
	out.SetFieldByName("entity", entity)

	return out, nil
}

func (s *configstoreTypedService) TypedDelete(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	messageFactory := dynamic.NewMessageFactoryWithDefaults()

	requestMessageDescriptor := s.GenResult.MessageMap[fmt.Sprintf("Delete%sRequest", s.KindName)]
	in := messageFactory.NewDynamicMessage(requestMessageDescriptor)
	if err := dec(in); err != nil {
		return nil, err
	}

	key, err := in.TryGetFieldByName("key")
	if err != nil {
		return nil, err
	}

	keyV, ok := key.(*dynamic.Message)
	if !ok {
		return nil, fmt.Errorf("unable to read key")
	}

	ref, err := convertKeyToDocumentRef(
		client,
		keyV,
	)
	if err != nil {
		return nil, err
	}

	// TODO: Validate ref is of the correct kind

	snapshot, err := ref.Get(ctx)
	if err != nil {
		return nil, err
	}

	entity, err := convertSnapshotToDynamicMessage(
		messageFactory,
		s.GenResult.MessageMap[s.KindName],
		snapshot,
		s.GenResult.CommonMessageDescriptors,
	)
	if err != nil {
		return nil, err
	}

	_, err = ref.Delete(ctx)
	if err != nil {
		return nil, err
	}

	responseMessageDescriptor := s.GenResult.MessageMap[fmt.Sprintf("Delete%sResponse", s.KindName)]
	out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
	out.SetFieldByName("entity", entity)

	return out, nil
}

func (s *configstoreTypedService) TypedWatch(srv interface{}, stream grpc.ServerStream) error {
	messageFactory := dynamic.NewMessageFactoryWithDefaults()

	snapshots := client.Collection(s.KindName).Snapshots(ctx)
	for true {
		snapshot, err := snapshots.Next()
		if err != nil {
			return err
		}
		for _, change := range snapshot.Changes {
			entity, err := convertSnapshotToDynamicMessage(
				messageFactory,
				s.GenResult.MessageMap[s.KindName],
				change.Doc,
				s.GenResult.CommonMessageDescriptors,
			)
			if err != nil {
				return err
			}

			responseMessageDescriptor := s.GenResult.MessageMap[fmt.Sprintf("Watch%sEvent", s.KindName)]
			out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
			switch change.Kind {
			case firestore.DocumentAdded:
				out.SetFieldByName("type", s.GenResult.WatchTypeEnumValues.Created.GetNumber())
			case firestore.DocumentModified:
				out.SetFieldByName("type", s.GenResult.WatchTypeEnumValues.Updated.GetNumber())
			case firestore.DocumentRemoved:
				out.SetFieldByName("type", s.GenResult.WatchTypeEnumValues.Deleted.GetNumber())
			}
			out.SetFieldByName("entity", entity)
			out.SetFieldByName("oldIndex", change.OldIndex)
			out.SetFieldByName("newIndex", change.NewIndex)

			stream.SendMsg(out)
		}
	}

	return nil
}
