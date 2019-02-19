package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"

	"github.com/gorilla/mux"

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

var client *firestore.Client

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
	genResult, err := generate(config.SchemaPath)
	if err != nil {
		log.Fatalln(err)
	}

	// Emit the testclient protobuf specification
	printer := new(protoprint.Printer)
	clientProtoFile, err := printer.PrintProtoToString(genResult.FileDesc)
	if err != nil {
		log.Fatalln(fmt.Sprintf("can't generate protobuf spec: %s", err))
	}
	clientProtoGoCode, err := generateGoCode(genResult.FileDesc, genResult.Schema)
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
		client, err = app.Firestore(ctx)
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
		for _, service := range genResult.Services {
			// kindSchema := kindMap[service]
			kindName := genResult.KindNameMap[service]

			grpcServer.RegisterService(
				&grpc.ServiceDesc{
					ServiceName: fmt.Sprintf("%s.%s", genResult.Schema.Name, service.GetName()),
					HandlerType: (*emptyServerInterface)(nil),
					Methods: []grpc.MethodDesc{
						{
							MethodName: "List",
							Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
								messageFactory := dynamic.NewMessageFactoryWithDefaults()

								requestMessageDescriptor := genResult.MessageMap[fmt.Sprintf("List%sRequest", kindName)]
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
									snapshots, err = client.Collection(kindName).Documents(ctx).GetAll()
								} else if limit == nil || limit.(uint32) == 0 {
									snapshots, err = client.Collection(kindName).OrderBy(firestore.DocumentID, firestore.Asc).StartAfter(start.(string)).Documents(ctx).GetAll()
								} else if start == nil {
									snapshots, err = client.Collection(kindName).Limit(int(limit.(uint32))).Documents(ctx).GetAll()
								} else {
									snapshots, err = client.Collection(kindName).OrderBy(firestore.DocumentID, firestore.Asc).StartAfter(start.(string)).Limit(int(limit.(uint32))).Documents(ctx).GetAll()
								}

								if err != nil {
									return nil, err
								}

								var entities []*dynamic.Message
								for _, snapshot := range snapshots {
									entity, err := convertSnapshotToDynamicMessage(
										messageFactory,
										genResult.MessageMap[kindName],
										snapshot,
										genResult.CommonMessageDescriptors,
									)
									if err != nil {
										return nil, err
									}
									entities = append(entities, entity)
								}

								responseMessageDescriptor := genResult.MessageMap[fmt.Sprintf("List%sResponse", kindName)]
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
							},
						},
						{
							MethodName: "Get",
							Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
								messageFactory := dynamic.NewMessageFactoryWithDefaults()

								requestMessageDescriptor := genResult.MessageMap[fmt.Sprintf("Get%sRequest", kindName)]
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
									genResult.MessageMap[kindName],
									snapshot,
									genResult.CommonMessageDescriptors,
								)
								if err != nil {
									return nil, err
								}

								responseMessageDescriptor := genResult.MessageMap[fmt.Sprintf("Get%sResponse", kindName)]
								out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
								out.SetFieldByName("entity", entity)

								return out, nil
							},
						},
						{
							MethodName: "Update",
							Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
								messageFactory := dynamic.NewMessageFactoryWithDefaults()

								requestMessageDescriptor := genResult.MessageMap[fmt.Sprintf("Update%sRequest", kindName)]
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
									genResult.MessageMap[kindName],
									entity.(*dynamic.Message),
									snapshot,
									genResult.Schema.Kinds[kindName],
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

								responseMessageDescriptor := genResult.MessageMap[fmt.Sprintf("Update%sResponse", kindName)]
								out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
								out.SetFieldByName("entity", entity)

								return out, nil
							},
						},
						{
							MethodName: "Create",
							Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
								messageFactory := dynamic.NewMessageFactoryWithDefaults()

								requestMessageDescriptor := genResult.MessageMap[fmt.Sprintf("Create%sRequest", kindName)]
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
									genResult.MessageMap[kindName],
									entity.(*dynamic.Message),
									nil,
									genResult.Schema.Kinds[kindName],
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
									genResult.CommonMessageDescriptors,
								)
								if err != nil {
									return nil, err
								}

								// set the ID back
								entity.(*dynamic.Message).SetFieldByName("key", key)

								responseMessageDescriptor := genResult.MessageMap[fmt.Sprintf("Create%sResponse", kindName)]
								out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
								out.SetFieldByName("entity", entity)

								return out, nil
							},
						},
						{
							MethodName: "Delete",
							Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
								messageFactory := dynamic.NewMessageFactoryWithDefaults()

								requestMessageDescriptor := genResult.MessageMap[fmt.Sprintf("Delete%sRequest", kindName)]
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
									genResult.MessageMap[kindName],
									snapshot,
									genResult.CommonMessageDescriptors,
								)
								if err != nil {
									return nil, err
								}

								_, err = ref.Delete(ctx)
								if err != nil {
									return nil, err
								}

								responseMessageDescriptor := genResult.MessageMap[fmt.Sprintf("Delete%sResponse", kindName)]
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
											genResult.MessageMap[kindName],
											change.Doc,
											genResult.CommonMessageDescriptors,
										)
										if err != nil {
											return err
										}

										responseMessageDescriptor := genResult.MessageMap[fmt.Sprintf("Watch%sEvent", kindName)]
										out := messageFactory.NewDynamicMessage(responseMessageDescriptor)
										switch change.Kind {
										case firestore.DocumentAdded:
											out.SetFieldByName("type", genResult.WatchTypeEnumValues.Created.GetNumber())
										case firestore.DocumentModified:
											out.SetFieldByName("type", genResult.WatchTypeEnumValues.Updated.GetNumber())
										case firestore.DocumentRemoved:
											out.SetFieldByName("type", genResult.WatchTypeEnumValues.Deleted.GetNumber())
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
					Metadata: genResult.FileBuilder.GetName(),
				},
				emptyServer,
			)
		}

		RegisterConfigstoreMetaServiceServer(grpcServer, &configstoreMetaServiceServer{
			schema: genResult.Schema,
		})

		// Start gRPC server.
		go func() {
			fmt.Println(fmt.Sprintf("Running gRPC server on port %d...", config.GrpcPort))
			grpcServer.Serve(lis)
		}()

		// Start HTTP server.
		r := mux.NewRouter()
		r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// We don't use http.ServeFile here, because it tries to take the current URL into
			// account when locating the file to serve. Since this is the 404 Not Found handler,
			// the current URL could be anything. We just always want to serve the contents
			// of index-react.html if this function is handling a request.
			f, err := os.Open("/server-ui/index.html")
			if err != nil {
				http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
				return
			}
			defer f.Close()

			d, err := f.Stat()
			if err != nil {
				http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
				return
			}

			http.ServeContent(w, r, d.Name(), d.ModTime(), f)
		})
		r.HandleFunc("/sdk/client.proto", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", clientProtoFile)
		})
		r.HandleFunc("/sdk/client.go", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", clientProtoGoCode)
		})
		http.Handle("/", r)
		fmt.Println(fmt.Sprintf("Running HTTP server on port %d...", config.HTTPPort))
		wrappedGrpc := grpcweb.WrapServer(grpcServer)
		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.HTTPPort), http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			resp.Header().Set("Access-Control-Allow-Origin", "*")
			resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, x-grpc-web")
			if req.Method == "OPTIONS" {
				return
			}
			if wrappedGrpc.IsGrpcWebRequest(req) {
				wrappedGrpc.ServeHTTP(resp, req)
				return
			}
			http.DefaultServeMux.ServeHTTP(resp, req)
		}))
	} else if mode == runModeGenerate {
		fmt.Println(clientProtoGoCode)
	}
}
