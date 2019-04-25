package main

import (
	"context"
	"fmt"

	"github.com/jhump/protoreflect/desc/builder"
	"github.com/jhump/protoreflect/dynamic"

	"google.golang.org/grpc"

	"cloud.google.com/go/firestore"
)

type configstoreDynamicProtobufService struct {
	firestoreClient      *firestore.Client
	genResult            *generatorResult
	service              *builder.ServiceBuilder
	kindName             string
	schema               *Schema
	transactionProcessor *transactionProcessor
	transactionWatcher   *transactionWatcher
}

func createConfigstoreDynamicProtobufServer(
	firestoreClient *firestore.Client,
	genResult *generatorResult,
	service *builder.ServiceBuilder,
	kindName string,
	schema *Schema,
	transactionWatcher *transactionWatcher,
) *configstoreDynamicProtobufService {
	return &configstoreDynamicProtobufService{
		firestoreClient:      firestoreClient,
		genResult:            genResult,
		service:              service,
		kindName:             kindName,
		schema:               schema,
		transactionProcessor: createTransactionProcessor(firestoreClient),
		transactionWatcher:   transactionWatcher,
	}
}

func (s *configstoreDynamicProtobufService) getMetaServiceServer() *configstoreMetaServiceServer {
	return createConfigstoreMetaServiceServer(
		s.firestoreClient,
		s.schema,
		s.transactionProcessor,
		s.transactionWatcher,
	)
}

func (s *configstoreDynamicProtobufService) dynamicProtobufList(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	messageFactory := dynamic.NewMessageFactoryWithDefaults()

	requestMessageDescriptor := s.genResult.MessageMap[fmt.Sprintf("List%sRequest", s.kindName)]
	in := messageFactory.NewDynamicMessage(requestMessageDescriptor)
	if err := dec(in); err != nil {
		return nil, err
	}

	startBytes, err := in.TryGetFieldByName("start")
	if err != nil {
		return nil, err
	}
	limitRaw, err := in.TryGetFieldByName("limit")
	if err != nil {
		return nil, err
	}

	var start []byte
	if startBytes != nil {
		start = startBytes.([]byte)
	}
	var limit uint32 = 0
	if limitRaw != nil {
		limit = limitRaw.(uint32)
	}

	metaServer := s.getMetaServiceServer()
	resp, err := metaServer.MetaList(ctx, &MetaListEntitiesRequest{
		Start:    start,
		Limit:    limit,
		KindName: s.kindName,
	})
	if err != nil {
		return nil, err
	}

	var entities []*dynamic.Message
	for _, snapshot := range resp.Entities {
		entity, err := convertMetaEntityToDynamicMessage(
			messageFactory,
			s.genResult.MessageMap[s.kindName],
			snapshot,
			s.genResult.CommonMessageDescriptors,
			s.genResult.KindMap[s.service],
		)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	responseMessageDescriptor := s.genResult.MessageMap[fmt.Sprintf("List%sResponse", s.kindName)]
	out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
	out.SetFieldByName("entities", entities)

	if !(limit <= 0) {
		if int64(len(entities)) < int64(limit) {
			out.SetFieldByName("moreResults", false)
		} else {
			// TODO: query to see if there really are more results, to make this behave like datastore
			out.SetFieldByName("moreResults", true)
			last := resp.Entities[len(resp.Entities)-1]
			out.SetFieldByName("next", []byte(last.Key.Path[len(last.Key.Path)-1].GetName()))
		}
	} else {
		out.SetFieldByName("moreResults", false)
	}

	return out, nil
}

func (s *configstoreDynamicProtobufService) dynamicProtobufGet(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	messageFactory := dynamic.NewMessageFactoryWithDefaults()

	requestMessageDescriptor := s.genResult.MessageMap[fmt.Sprintf("Get%sRequest", s.kindName)]
	in := messageFactory.NewDynamicMessage(requestMessageDescriptor)
	if err := dec(in); err != nil {
		return nil, err
	}

	rawKey, err := in.TryGetFieldByName("key")
	if err != nil {
		return nil, err
	}

	key, ok := rawKey.(*Key)
	if !ok {
		return nil, fmt.Errorf("unable to read key")
	}

	metaServer := s.getMetaServiceServer()
	resp, err := metaServer.MetaGet(ctx, &MetaGetEntityRequest{
		Key:      key,
		KindName: s.kindName,
	})
	if err != nil {
		return nil, err
	}

	entity, err := convertMetaEntityToDynamicMessage(
		messageFactory,
		s.genResult.MessageMap[s.kindName],
		resp.Entity,
		s.genResult.CommonMessageDescriptors,
		s.genResult.KindMap[s.service],
	)
	if err != nil {
		return nil, err
	}

	responseMessageDescriptor := s.genResult.MessageMap[fmt.Sprintf("Get%sResponse", s.kindName)]
	out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
	out.SetFieldByName("entity", entity)

	return out, nil
}

func (s *configstoreDynamicProtobufService) dynamicProtobufUpdate(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	messageFactory := dynamic.NewMessageFactoryWithDefaults()

	requestMessageDescriptor := s.genResult.MessageMap[fmt.Sprintf("Update%sRequest", s.kindName)]
	in := messageFactory.NewDynamicMessage(requestMessageDescriptor)
	if err := dec(in); err != nil {
		return nil, err
	}

	rawEntity, err := in.TryGetFieldByName("entity")
	if err != nil {
		return nil, err
	}

	if rawEntity == nil {
		return nil, fmt.Errorf("entity must not be nil")
	}

	entity, err := convertDynamicMessageIntoMetaEntity(
		s.firestoreClient,
		messageFactory,
		s.genResult.MessageMap[s.kindName],
		rawEntity.(*dynamic.Message),
		s.genResult.Schema.Kinds[s.kindName],
	)
	if err != nil {
		return nil, err
	}

	metaServer := s.getMetaServiceServer()
	resp, err := metaServer.MetaUpdate(ctx, &MetaUpdateEntityRequest{
		Entity: entity,
	})
	if err != nil {
		return nil, err
	}

	message, err := convertMetaEntityToDynamicMessage(
		messageFactory,
		s.genResult.MessageMap[s.kindName],
		resp.Entity,
		s.genResult.CommonMessageDescriptors,
		s.genResult.KindMap[s.service],
	)
	if err != nil {
		return nil, err
	}

	responseMessageDescriptor := s.genResult.MessageMap[fmt.Sprintf("Update%sResponse", s.kindName)]
	out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
	out.SetFieldByName("entity", message)

	return out, nil
}

func (s *configstoreDynamicProtobufService) dynamicProtobufCreate(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	messageFactory := dynamic.NewMessageFactoryWithDefaults()

	requestMessageDescriptor := s.genResult.MessageMap[fmt.Sprintf("Create%sRequest", s.kindName)]
	in := messageFactory.NewDynamicMessage(requestMessageDescriptor)
	if err := dec(in); err != nil {
		return nil, err
	}

	rawEntity, err := in.TryGetFieldByName("entity")
	if err != nil {
		return nil, err
	}

	if rawEntity == nil {
		return nil, fmt.Errorf("entity must not be nil")
	}

	entity, err := convertDynamicMessageIntoMetaEntity(
		s.firestoreClient,
		messageFactory,
		s.genResult.MessageMap[s.kindName],
		rawEntity.(*dynamic.Message),
		s.genResult.Schema.Kinds[s.kindName],
	)
	if err != nil {
		return nil, err
	}

	metaServer := s.getMetaServiceServer()
	resp, err := metaServer.MetaCreate(ctx, &MetaCreateEntityRequest{
		Entity:   entity,
		KindName: s.kindName,
	})
	if err != nil {
		return nil, err
	}

	message, err := convertMetaEntityToDynamicMessage(
		messageFactory,
		s.genResult.MessageMap[s.kindName],
		resp.Entity,
		s.genResult.CommonMessageDescriptors,
		s.genResult.KindMap[s.service],
	)
	if err != nil {
		return nil, err
	}

	responseMessageDescriptor := s.genResult.MessageMap[fmt.Sprintf("Create%sResponse", s.kindName)]
	out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
	out.SetFieldByName("entity", message)

	return out, nil
}

func (s *configstoreDynamicProtobufService) dynamicProtobufDelete(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	messageFactory := dynamic.NewMessageFactoryWithDefaults()

	requestMessageDescriptor := s.genResult.MessageMap[fmt.Sprintf("Delete%sRequest", s.kindName)]
	in := messageFactory.NewDynamicMessage(requestMessageDescriptor)
	if err := dec(in); err != nil {
		return nil, err
	}

	rawKey, err := in.TryGetFieldByName("key")
	if err != nil {
		return nil, err
	}

	key, ok := rawKey.(*Key)
	if !ok {
		return nil, fmt.Errorf("unable to read key")
	}

	metaServer := s.getMetaServiceServer()
	resp, err := metaServer.MetaDelete(ctx, &MetaDeleteEntityRequest{
		Key:      key,
		KindName: s.kindName,
	})
	if err != nil {
		return nil, err
	}

	entity, err := convertMetaEntityToDynamicMessage(
		messageFactory,
		s.genResult.MessageMap[s.kindName],
		resp.Entity,
		s.genResult.CommonMessageDescriptors,
		s.genResult.KindMap[s.service],
	)
	if err != nil {
		return nil, err
	}

	responseMessageDescriptor := s.genResult.MessageMap[fmt.Sprintf("Delete%sResponse", s.kindName)]
	out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
	out.SetFieldByName("entity", entity)

	return out, nil
}

func (s *configstoreDynamicProtobufService) dynamicProtobufWatch(srv interface{}, ctx context.Context, stream grpc.ServerStream) error {
	messageFactory := dynamic.NewMessageFactoryWithDefaults()

	snapshots := s.firestoreClient.Collection(s.kindName).Snapshots(ctx)
	for true {
		snapshot, err := snapshots.Next()
		if err != nil {
			return err
		}
		for _, change := range snapshot.Changes {
			metaEntity, err := convertSnapshotToMetaEntity(
				s.genResult.KindMap[s.service],
				change.Doc,
			)
			if err != nil {
				return err
			}
			message, err := convertMetaEntityToDynamicMessage(
				messageFactory,
				s.genResult.MessageMap[s.kindName],
				metaEntity,
				s.genResult.CommonMessageDescriptors,
				s.genResult.KindMap[s.service],
			)
			if err != nil {
				return err
			}

			responseMessageDescriptor := s.genResult.MessageMap[fmt.Sprintf("Watch%sEvent", s.kindName)]
			out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
			switch change.Kind {
			case firestore.DocumentAdded:
				out.SetFieldByName("type", s.genResult.WatchTypeEnumValues.Created.GetNumber())
			case firestore.DocumentModified:
				out.SetFieldByName("type", s.genResult.WatchTypeEnumValues.Updated.GetNumber())
			case firestore.DocumentRemoved:
				out.SetFieldByName("type", s.genResult.WatchTypeEnumValues.Deleted.GetNumber())
			}
			out.SetFieldByName("entity", message)
			out.SetFieldByName("oldIndex", change.OldIndex)
			out.SetFieldByName("newIndex", change.NewIndex)

			stream.SendMsg(out)
		}
	}

	return nil
}
