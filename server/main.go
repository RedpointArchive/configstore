package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"cloud.google.com/go/firestore"

	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"

	firebase "firebase.google.com/go"
)

type emptyServerInterface interface {
}

type runtimeConfig struct {
	GoogleCloudProjectID          string `envconfig:"GOOGLE_CLOUD_PROJECT_ID" required:"true"`
	GoogleCloudServiceAccountPath string `envconfig:"GOOGLE_CLOUD_SERVICE_ACCOUNT_PATH"`
	GrpcPort                      uint16 `envconfig:"GRPC_PORT" required:"true"`
	HTTPPort                      uint16 `envconfig:"HTTP_PORT" required:"true"`
	SchemaPath                    string `envconfig:"SCHEMA_PATH" required:"true"`
}

type runMode string

const (
	runModeServe    runMode = "serve"
	runModeGenerate runMode = "generate"
)

func main() {
	mode := runModeServe
	generateFlag := flag.Bool("generate", false, "emit Go client code instead of serving traffic")
	flag.Parse()
	if *generateFlag {
		mode = runModeGenerate
	}

	config := &runtimeConfig{}
	err := envconfig.Process("CONFIGSTORE", config)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	// Generate the schema and gRPC types based on schema.json
	_, services, fileBuilder, fileDesc, schema, _, kindNameMap, messageMap, watchTypeEnumValues, err := generate(config.SchemaPath)
	if err != nil {
		log.Fatalln(err)
	}

	// Emit the testclient protobuf specification
	printer := new(protoprint.Printer)
	clientProtoFile, err := printer.PrintProtoToString(fileDesc)
	if err != nil {
		log.Fatalln(err)
	}
	clientProtoGoCode, err := generateGoCode(fileDesc, schema.Name)
	if err != nil {
		log.Fatalln(fmt.Sprintf("can't generate Go code: %s", err))
	}

	if mode == runModeServe {
		// Connect to Firestore using Application Default Credentials
		conf := &firebase.Config{ProjectID: config.GoogleCloudProjectID}
		app, err := firebase.NewApp(ctx, conf)
		if err != nil {
			log.Fatalln(err)
		}
		client, err := app.Firestore(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		defer client.Close()

		// Serve the configstore gRPC server
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GrpcPort))
		if err != nil {
			log.Fatalln(err)
		}
		grpcServer := grpc.NewServer()
		emptyServer := new(emptyServerInterface)
		for _, service := range services {
			// kindSchema := kindMap[service]
			kindName := kindNameMap[service]

			grpcServer.RegisterService(
				&grpc.ServiceDesc{
					ServiceName: fmt.Sprintf("%s.%s", schema.Name, service.GetName()),
					HandlerType: (*emptyServerInterface)(nil),
					Methods: []grpc.MethodDesc{
						{
							MethodName: "List",
							Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
								messageFactory := dynamic.NewMessageFactoryWithDefaults()

								requestMessageDescriptor := messageMap[fmt.Sprintf("List%sRequest", kindName)]
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

								if limit.(uint32) == 0 {
									return nil, fmt.Errorf("limit must be greater than 0, or omitted")
								}

								var start interface{}
								if startBytes != nil {
									if len(startBytes.([]byte)[:]) > 0 {
										start = string(startBytes.([]byte)[:])
									}
								}

								var snapshots []*firestore.DocumentSnapshot
								if limit == nil && start == nil {
									snapshots, err = client.Collection(kindName).Documents(ctx).GetAll()
								} else if limit == nil {
									ref := client.Doc(start.(string))
									snapshots, err = client.Collection(kindName).OrderBy(firestore.DocumentID, firestore.Asc).StartAfter(ref.ID).Documents(ctx).GetAll()
								} else if start == nil {
									snapshots, err = client.Collection(kindName).Limit(int(limit.(uint32))).Documents(ctx).GetAll()
								} else {
									ref := client.Doc(start.(string))
									snapshots, err = client.Collection(kindName).OrderBy(firestore.DocumentID, firestore.Asc).StartAfter(ref.ID).Limit(int(limit.(uint32))).Documents(ctx).GetAll()
								}

								if err != nil {
									return nil, err
								}

								var entities []*dynamic.Message
								for _, snapshot := range snapshots {
									entity, err := convertSnapshotToDynamicMessage(
										messageFactory,
										messageMap[kindName],
										snapshot,
									)
									if err != nil {
										return nil, err
									}
									entities = append(entities, entity)
								}

								responseMessageDescriptor := messageMap[fmt.Sprintf("List%sResponse", kindName)]
								out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
								out.SetFieldByName("entities", entities)

								if limit != nil {
									if uint32(len(entities)) < limit.(uint32) {
										out.SetFieldByName("moreResults", false)
									} else {
										// TODO: query to see if there really are more results, to make this behave like datastore
										out.SetFieldByName("moreResults", true)
										last := snapshots[len(snapshots)-1]
										out.SetFieldByName("next", []byte(last.Ref.Path))
									}
								} else {
									out.SetFieldByName("moreResults", false)
								}

								return out, nil
							},
						},
						{
							MethodName: "Get",
							Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
								messageFactory := dynamic.NewMessageFactoryWithDefaults()

								requestMessageDescriptor := messageMap[fmt.Sprintf("Get%sRequest", kindName)]
								in := messageFactory.NewDynamicMessage(requestMessageDescriptor)
								if err := dec(in); err != nil {
									return nil, err
								}

								id, err := in.TryGetFieldByName("id")
								if err != nil {
									return nil, err
								}

								snapshot, err := client.Collection(kindName).Doc(id.(string)).Get(ctx)
								if err != nil {
									return nil, err
								}

								entity, err := convertSnapshotToDynamicMessage(
									messageFactory,
									messageMap[kindName],
									snapshot,
								)
								if err != nil {
									return nil, err
								}

								responseMessageDescriptor := messageMap[fmt.Sprintf("Get%sResponse", kindName)]
								out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
								out.SetFieldByName("entity", entity)

								return out, nil
							},
						},
						{
							MethodName: "Update",
							Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
								messageFactory := dynamic.NewMessageFactoryWithDefaults()

								requestMessageDescriptor := messageMap[fmt.Sprintf("Update%sRequest", kindName)]
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

								id, data, err := convertDynamicMessageIntoIDAndDataMap(
									messageFactory,
									messageMap[kindName],
									entity.(*dynamic.Message),
								)

								if id == "" {
									return nil, fmt.Errorf("entity must be set")
								}

								_, err = client.Collection(kindName).Doc(id).Set(ctx, data)
								if err != nil {
									return nil, err
								}

								responseMessageDescriptor := messageMap[fmt.Sprintf("Update%sResponse", kindName)]
								out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
								out.SetFieldByName("entity", entity)

								return out, nil
							},
						},
						{
							MethodName: "Create",
							Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
								messageFactory := dynamic.NewMessageFactoryWithDefaults()

								requestMessageDescriptor := messageMap[fmt.Sprintf("Create%sRequest", kindName)]
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

								id, data, err := convertDynamicMessageIntoIDAndDataMap(
									messageFactory,
									messageMap[kindName],
									entity.(*dynamic.Message),
								)

								if id != "" {
									return nil, fmt.Errorf("entity must be nil / empty / unset")
								}

								ref, _, err := client.Collection(kindName).Add(ctx, data)
								if err != nil {
									return nil, err
								}

								// set the ID back
								entity.(*dynamic.Message).SetFieldByName("id", ref.ID)

								responseMessageDescriptor := messageMap[fmt.Sprintf("Create%sResponse", kindName)]
								out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
								out.SetFieldByName("entity", entity)

								return out, nil
							},
						},
						{
							MethodName: "Delete",
							Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
								messageFactory := dynamic.NewMessageFactoryWithDefaults()

								requestMessageDescriptor := messageMap[fmt.Sprintf("Delete%sRequest", kindName)]
								in := messageFactory.NewDynamicMessage(requestMessageDescriptor)
								if err := dec(in); err != nil {
									return nil, err
								}

								id, err := in.TryGetFieldByName("id")
								if err != nil {
									return nil, err
								}

								snapshot, err := client.Collection(kindName).Doc(id.(string)).Get(ctx)
								if err != nil {
									return nil, err
								}

								entity, err := convertSnapshotToDynamicMessage(
									messageFactory,
									messageMap[kindName],
									snapshot,
								)
								if err != nil {
									return nil, err
								}

								_, err = client.Collection(kindName).Doc(id.(string)).Delete(ctx)
								if err != nil {
									return nil, err
								}

								responseMessageDescriptor := messageMap[fmt.Sprintf("Delete%sResponse", kindName)]
								out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
								out.SetFieldByName("entity", entity)

								return out, nil
							},
						},
					},
					Streams: []grpc.StreamDesc{
						{
							StreamName:    "Watch",
							ServerStreams: true,
							ClientStreams: false,
							Handler: func(srv interface{}, stream grpc.ServerStream) error {
								messageFactory := dynamic.NewMessageFactoryWithDefaults()

								snapshots := client.Collection(kindName).Snapshots(ctx)
								for true {
									snapshot, err := snapshots.Next()
									if err != nil {
										return err
									}
									for _, change := range snapshot.Changes {
										entity, err := convertSnapshotToDynamicMessage(
											messageFactory,
											messageMap[kindName],
											change.Doc,
										)
										if err != nil {
											return err
										}

										responseMessageDescriptor := messageMap[fmt.Sprintf("Watch%sEvent", kindName)]
										out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
										switch change.Kind {
										case firestore.DocumentAdded:
											out.SetFieldByName("type", watchTypeEnumValues.Created.GetNumber())
										case firestore.DocumentModified:
											out.SetFieldByName("type", watchTypeEnumValues.Updated.GetNumber())
										case firestore.DocumentRemoved:
											out.SetFieldByName("type", watchTypeEnumValues.Deleted.GetNumber())
										}
										out.SetFieldByName("entity", entity)
										out.SetFieldByName("oldIndex", change.OldIndex)
										out.SetFieldByName("newIndex", change.NewIndex)

										stream.SendMsg(out)
									}
								}

								return nil
							},
						},
					},
					Metadata: fileBuilder.GetName(),
				},
				emptyServer,
			)
		}

		// Start gRPC server.
		go func() {
			fmt.Println(fmt.Sprintf("Running gRPC server on port %d...", config.GrpcPort))
			grpcServer.Serve(lis)
		}()

		// Start HTTP server.
		http.HandleFunc("/sdk/client.proto", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", clientProtoFile)
		})
		http.HandleFunc("/sdk/client.go", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", clientProtoGoCode)
		})
		fmt.Println(fmt.Sprintf("Running HTTP server on port %d...", config.HTTPPort))
		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.HTTPPort), nil)
	} else if mode == runModeGenerate {
		fmt.Println(clientProtoGoCode)
	}
}
