package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"

	"github.com/gorilla/handlers"
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

func main() {
	mode := runModeServe
	generateFlag := flag.Bool("generate", false, "emit Go client code instead of serving traffic")
	reverseProxyDevServerFlag := flag.Bool("enable-reverse-proxy-react-dev-server", false, "instead of serving React files from /server-ui/, reverse proxy to http://configstore-dashboard:3000/")
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
		for _, service := range genResult.Services {
			dynamicProtobufServer := createConfigstoreDynamicProtobufServer(
				client,
				genResult,
				service,
				genResult.KindNameMap[service],
				genResult.Schema,
			)

			grpcServer.RegisterService(
				&grpc.ServiceDesc{
					ServiceName: fmt.Sprintf("%s.%s", genResult.Schema.Name, service.GetName()),
					HandlerType: (*emptyServerInterface)(nil),
					Methods: []grpc.MethodDesc{
						{
							MethodName: "List",
							Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
								out, err := dynamicProtobufServer.dynamicProtobufList(srv, ctx, dec, interceptor)
								return out, err
							},
						},
						{
							MethodName: "Get",
							Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
								out, err := dynamicProtobufServer.dynamicProtobufGet(srv, ctx, dec, interceptor)
								return out, err
							},
						},
						{
							MethodName: "Update",
							Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
								out, err := dynamicProtobufServer.dynamicProtobufUpdate(srv, ctx, dec, interceptor)
								return out, err
							},
						},
						{
							MethodName: "Create",
							Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
								out, err := dynamicProtobufServer.dynamicProtobufCreate(srv, ctx, dec, interceptor)
								return out, err
							},
						},
						{
							MethodName: "Delete",
							Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
								out, err := dynamicProtobufServer.dynamicProtobufDelete(srv, ctx, dec, interceptor)
								return out, err
							},
						},
					},
					Streams: []grpc.StreamDesc{
						{
							StreamName:    "Watch",
							ServerStreams: true,
							ClientStreams: false,
							Handler: func(srv interface{}, stream grpc.ServerStream) error {
								return dynamicProtobufServer.dynamicProtobufWatch(srv, ctx, stream)
							},
						},
					},
					Metadata: genResult.FileBuilder.GetName(),
				},
				emptyServer,
			)
		}

		// Add the metadata server.
		RegisterConfigstoreMetaServiceServer(grpcServer, createConfigstoreMetaServiceServer(
			client,
			genResult.Schema,
			createTransactionProcessor(client),
		))

		// Start gRPC server.
		go func() {
			fmt.Println(fmt.Sprintf("Running gRPC server on port %d...", config.GrpcPort))
			grpcServer.Serve(lis)
		}()

		// Start HTTP server.
		router := mux.NewRouter()
		/*r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// We don't use http.ServeFile here, because it tries to take the current URL into
			// account when locating the file to serve. Since this is the 404 Not Found handler,
			// the current URL could be anything. We just always want to serve the contents
			// of index.html if this function is handling a request.
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
		})*/

		router.HandleFunc("/sdk/client.proto", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", clientProtoFile)
		})
		router.HandleFunc("/sdk/client.go", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", clientProtoGoCode)
		})

		if *reverseProxyDevServerFlag {
			fmt.Println("Enabling reverse proxy to React dev server at http://configstore-dashboard:3000/...")
			url, _ := url.Parse("http://configstore-dashboard:3000/")
			rp := httputil.NewSingleHostReverseProxy(url)
			router.PathPrefix("/static").Handler(rp)
			router.PathPrefix("/sockjs-node").Handler(rp)
			router.PathPrefix("/").Handler(rp)
		} else {
			router.PathPrefix("/static").Handler(http.FileServer(http.Dir("/server-ui/")))
			router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, "/server-ui/index.html")
			})
		}

		wrappedGrpc := grpcweb.WrapServer(
			grpcServer,
			grpcweb.WithAllowedRequestHeaders([]string{"*"}),
			grpcweb.WithCorsForRegisteredEndpointsOnly(false),
			grpcweb.WithOriginFunc(func(origin string) bool {
				return true
			}),
			grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
				return true
			}),
			grpcweb.WithWebsockets(true),
		)
		root := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			if wrappedGrpc.IsGrpcWebRequest(req) ||
				wrappedGrpc.IsAcceptableGrpcCorsRequest(req) ||
				wrappedGrpc.IsGrpcWebSocketRequest(req) {
				resp.Header().Set("Content-Type", "application/grpc-web-text")
				wrappedGrpc.ServeHTTP(resp, req)
				return
			}
			router.ServeHTTP(resp, req)
		})

		fmt.Println(fmt.Sprintf("Running HTTP server on port %d...", config.HTTPPort))
		srv := &http.Server{
			Handler: handlers.LoggingHandler(os.Stdout, root),
			Addr:    fmt.Sprintf("0.0.0.0:%d", config.HTTPPort),
		}

		err = srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		} else {
			log.Print("HTTP server shutdown gracefully.")
		}
	} else if mode == runModeGenerate {
		fmt.Println(clientProtoGoCode)
	}
}
