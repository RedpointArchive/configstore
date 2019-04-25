package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dgrijalva/jwt-go"
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
	AllowedOrigins                string `envconfig:"ALLOWED_ORIGINS"`
	AuthEnabled                   bool   `envconfig:"AUTH_ENABLED"`
	AuthSigningKey                string `envconfig:"AUTH_SIGNING_KEY"`
	AuthRequiredScopes            string `envconfig:"AUTH_REQUIRED_SCOPES"`
	AuthIss                       string `envconfig:"AUTH_ISS"`
	AuthAud                       string `envconfig:"AUTH_AUD"`
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
		log.Fatalln(fmt.Errorf("can't read configuration from environment: %v", err))
	}

	ctx := context.Background()

	// HACK: Workaround Auth0 timestamp issues
	jwt.TimeFunc = func() time.Time {
		// Pretend we're 1 minute in the future because Auth0 can sometimes issue
		// tokens in the future, and that prevents us from ever authenticating correctly,
		// even if the token is only like 2 seconds in the future.
		return time.Now().Add(time.Minute * 1)
	}

	// Generate the schema and gRPC types based on schema.json
	genResult, err := generate(config.SchemaPath)
	if err != nil {
		log.Fatalln(fmt.Errorf("can't generate protobufs: %v", err))
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
			log.Fatalln(fmt.Errorf("can't connect to Firebase: %v", err))
		}
		client, err := app.Firestore(ctx)
		if err != nil {
			log.Fatalln(fmt.Errorf("can't connect to Firestore: %v", err))
		}
		defer client.Close()

		// Start the transaction watcher
		transactionWatcher, err := createTransactionWatcher(ctx, client, genResult.Schema)
		if err != nil {
			log.Fatalln(fmt.Errorf("can't create transaction watcher: %v", err))
		}

		// Serve the configstore gRPC server
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GrpcPort))
		if err != nil {
			log.Fatalln(fmt.Errorf("can't serve the gRPC server: %v", err))
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
								return dynamicProtobufServer.dynamicProtobufWatch(srv, stream.Context(), stream)
							},
						},
					},
					Metadata: genResult.FileBuilder.GetName(),
				},
				emptyServer,
			)
		}

		dynamicProtobufTransactionServer := createConfigstoreDynamicProtobufTransactionServer(
			client,
			genResult,
			genResult.TransactionService,
			genResult.Schema,
			transactionWatcher,
		)
		grpcServer.RegisterService(
			&grpc.ServiceDesc{
				ServiceName: fmt.Sprintf("%s.%s", genResult.Schema.Name, genResult.TransactionService.GetName()),
				HandlerType: (*emptyServerInterface)(nil),
				Methods:     nil,
				Streams: []grpc.StreamDesc{
					{
						StreamName:    "Watch",
						ServerStreams: true,
						ClientStreams: false,
						Handler: func(srv interface{}, stream grpc.ServerStream) error {
							return dynamicProtobufTransactionServer.dynamicProtobufTransactionWatch(stream.Context(), srv, stream)
						},
					},
				},
				Metadata: genResult.FileBuilder.GetName(),
			},
			emptyServer,
		)

		// Add the metadata server.
		metaServer := createConfigstoreMetaServiceServer(
			client,
			genResult.Schema,
			createTransactionProcessor(client),
			transactionWatcher,
		)
		RegisterConfigstoreMetaServiceServer(grpcServer, metaServer)
		grpcServer.RegisterService(&grpc.ServiceDesc{
			ServiceName: fmt.Sprintf("%s.%s", genResult.Schema.Name, "ConfigstoreMetaService"),
			HandlerType: (*ConfigstoreMetaServiceServer)(nil),
			Methods: []grpc.MethodDesc{
				{
					MethodName: "GetSchema",
					Handler:    _ConfigstoreMetaService_GetSchema_Handler,
				},
				{
					MethodName: "MetaList",
					Handler:    _ConfigstoreMetaService_MetaList_Handler,
				},
				{
					MethodName: "MetaGet",
					Handler:    _ConfigstoreMetaService_MetaGet_Handler,
				},
				{
					MethodName: "MetaUpdate",
					Handler:    _ConfigstoreMetaService_MetaUpdate_Handler,
				},
				{
					MethodName: "MetaCreate",
					Handler:    _ConfigstoreMetaService_MetaCreate_Handler,
				},
				{
					MethodName: "MetaDelete",
					Handler:    _ConfigstoreMetaService_MetaDelete_Handler,
				},
				{
					MethodName: "GetDefaultPartitionId",
					Handler:    _ConfigstoreMetaService_GetDefaultPartitionId_Handler,
				},
				{
					MethodName: "ApplyTransaction",
					Handler:    _ConfigstoreMetaService_ApplyTransaction_Handler,
				},
			},
			Streams: []grpc.StreamDesc{
				{
					StreamName:    "WatchTransactions",
					Handler:       _ConfigstoreMetaService_WatchTransactions_Handler,
					ServerStreams: true,
				},
			},
			Metadata: "meta.proto",
		}, metaServer)

		// Start gRPC server.
		go func() {
			fmt.Println(fmt.Sprintf("Running gRPC server on port %d...", config.GrpcPort))
			grpcServer.Serve(lis)
		}()

		// Start HTTP server.
		router := mux.NewRouter()
		router.HandleFunc("/sdk/client.proto", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", clientProtoFile)
		})
		router.HandleFunc("/sdk/client.go", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", clientProtoGoCode)
		})

		router.PathPrefix("/static").Handler(http.FileServer(http.Dir("/server-ui/")))
		router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "/server-ui/index.html")
		})

		fmt.Println(fmt.Sprintf("Running HTTP server on port %d...", config.HTTPPort))

		GrpcServeWithWrapper(router, grpcServer, func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if config.AuthEnabled {
					authMiddleware(
						h,
						[]byte(config.AuthSigningKey),
						strings.Split(config.AuthRequiredScopes, ","),
						config.AuthIss,
						config.AuthAud,
					).ServeHTTP(w, r)
				} else {
					h.ServeHTTP(w, r)
				}
			})
		}, fmt.Sprintf("0.0.0.0:%d", config.HTTPPort), func(origin string) bool {
			if config.AllowedOrigins == "" {
				return false
			}
			if config.AllowedOrigins == "*" {
				return true
			}
			allowedOriginsArray := strings.Split(config.AllowedOrigins, ",")
			for _, or := range allowedOriginsArray {
				if origin == or {
					return true
				}
			}
			return false
		})
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

type wrapCall func(http.Handler) http.Handler

// GrpcServeWithWrapper allows you to wrap the gRPC layer in additional HTTP handler. For example,
// you can use this to inject an authentication layer for gRPC.
func GrpcServeWithWrapper(handler http.Handler, grpcServer *grpc.Server, wrapper wrapCall, addr string, isAllowedOrigin func(origin string) bool) {
	wrappedGrpc := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithAllowedRequestHeaders([]string{"*"}),
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(origin string) bool {
			return isAllowedOrigin(origin)
		}),
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
			origin, err := grpcweb.WebsocketRequestOrigin(req)
			if err != nil {
				return false
			}
			return isAllowedOrigin(origin)
		}),
		grpcweb.WithWebsockets(true),
	)
	doubleWrappedGrpc := wrapper(wrappedGrpc)
	mixHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		// Serve GRPC if this is a GRPC request.
		if wrappedGrpc.IsAcceptableGrpcCorsRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
			return
		}
		if wrappedGrpc.IsGrpcWebRequest(req) ||
			wrappedGrpc.IsGrpcWebSocketRequest(req) {
			setDummyClientHeaders(resp)
			doubleWrappedGrpc.ServeHTTP(resp, req)
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
