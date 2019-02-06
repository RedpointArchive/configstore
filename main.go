package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"

	firebase "firebase.google.com/go"
)

type emptyServerInterface interface {
}

func main() {
	ctx := context.Background()

	// Generate the schema and gRPC types based on schema.json
	_, services, fileBuilder, fileDesc, schema, _, kindNameMap, messageMap, err := generate()
	if err != nil {
		log.Fatalln(err)
	}

	// Emit the testclient protobuf specification
	printer := new(protoprint.Printer)
	protoFile, err := os.Create("testclient/testclient.proto")
	if err != nil {
		log.Fatalln(err)
	}
	defer protoFile.Close()
	fmt.Println("writing testclient/testclient.proto...")
	err = printer.PrintProtoFile(fileDesc, protoFile)
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to Firestore using Application Default Credentials
	conf := &firebase.Config{ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT_ID")}
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
	lis, err := net.Listen("tcp", ":13389")
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
							in := dynamic.NewMessage(messageMap[fmt.Sprintf("List%sRequest", kindName)])
							if err := dec(in); err != nil {
								return nil, err
							}
							fmt.Printf("%s", in.String())
							/*


								in := new(ListProjectAccessRequest)
								if err := dec(in); err != nil {
									return nil, err
								}
								if interceptor == nil {
									return srv.(ProjectAccessServiceServer).List(ctx, in)
								}
								info := &grpc.UnaryServerInfo{
									Server:     srv,
									FullMethod: "/configstoreExample.ProjectAccessService/List",
								}
								handler := func(ctx context.Context, req interface{}) (interface{}, error) {
									return srv.(ProjectAccessServiceServer).List(ctx, req.(*ListProjectAccessRequest))
								}
								return interceptor(ctx, in, info, handler)




								fmt.Println("%s", dec)*/

							return nil, fmt.Errorf("not implemented")
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
							return nil, fmt.Errorf("not implemented")
						},
					},
					{
						MethodName: "Create",
						Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
							return nil, fmt.Errorf("not implemented")
						},
					},
					{
						MethodName: "Delete",
						Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
							return nil, fmt.Errorf("not implemented")
						},
					},
				},
				Streams:  []grpc.StreamDesc{},
				Metadata: fileBuilder.GetName(),
			},
			emptyServer,
		)
	}
	fmt.Println("running server on port 13389...")
	grpcServer.Serve(lis)
}
