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
	"os/signal"
	"syscall"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"

	"time"

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

		// Start the transaction watcher
		transactionWatcher, err := createTransactionWatcher(ctx, client, genResult.Schema)
		if err != nil {
			log.Fatalln(err)
		}

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
				transactionWatcher,
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
			transactionWatcher,
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

		fmt.Println(fmt.Sprintf("Running HTTP server on port %d...", config.HTTPPort))

		GrpcServe(router, grpcServer, fmt.Sprintf("0.0.0.0:%d", config.HTTPPort))

	} else if mode == runModeGenerate {
		fmt.Println(clientProtoGoCode)
	}
}

// setDummyClientHeaders sets two empty headers "grpc-status" and "grpc-message", which
// the JS client tries to read on responses. We don't actually populate them with anything
// but the server-side gRPC handler looks at the headers that were sent by the server
// to determine what value to set in Access-Control-Expose-Headers. Since prior to this
// change we didn't send them at all, they wouldn't be included in Access-Control-Expose-Headers
// and thus when the JS client tried to read them it would log errors in the Chrome console
// (and those error counts add up pretty fast when we're repeatedly polling endpoints).
//
// These headers are apparently only used for streaming requests, and we don't use those
// yet (and ideally, they should get overwritten by the server-side gRPC if they need to be
// set).
func setDummyClientHeaders(resp http.ResponseWriter) {
	resp.Header().Set("grpc-status", "")
	resp.Header().Set("grpc-message", "")
}

// GrpcServe starts a HTTP server, with a gRPC server attached. This uses HttpServe
// underneath, and just contains the logic for setting up the HTTP handler with gRPC
// support.
func GrpcServe(handler http.Handler, grpcServer *grpc.Server, addr string) {
	wrappedGrpc := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithAllowedRequestHeaders([]string{"*"}),
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}),
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
			_, err := grpcweb.WebsocketRequestOrigin(req)
			if err != nil {
				return false
			}
			return true
		}),
		grpcweb.WithWebsockets(true),
	)
	mixHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		// Serve GRPC if this is a GRPC request.
		if wrappedGrpc.IsGrpcWebRequest(req) ||
			wrappedGrpc.IsAcceptableGrpcCorsRequest(req) ||
			wrappedGrpc.IsGrpcWebSocketRequest(req) {
			setDummyClientHeaders(resp)
			wrappedGrpc.ServeHTTP(resp, req)
			return
		}

		// Otherwise, fallback to the specified HTTP handler.
		handler.ServeHTTP(resp, req)
	})
	HttpServe(mixHandler, addr)
}

// HttpServe starts a HTTP server using the given HTTP handler and listening
// on the specified address.
func HttpServe(handler http.Handler, addr string) {
	server := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("http server error: %v", err)
			os.Exit(1)
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)

	<-signalChannel

	if os.Getenv("DEV") != "1" {
		fmt.Printf("gracefully shutting down...\n")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		if err.Error() != "http: Server closed" { // ignore this "error"
			fmt.Printf("could not gracefully shut down server: %v\n", err)
		}
	}
}
