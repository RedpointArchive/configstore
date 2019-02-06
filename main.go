package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/jhump/protoreflect/desc/builder"
	"github.com/jhump/protoreflect/desc/protoprint"
	"google.golang.org/grpc"
)

type emptyServerInterface interface {
}

func convertToType(t configstoreSchemaKindFieldType) *builder.FieldType {
	switch t {
	case typeDouble:
		return builder.FieldTypeDouble()
	case typeInt64:
		return builder.FieldTypeInt64()
	case typeString:
		return builder.FieldTypeString()
	case typeTimestamp:
		return builder.FieldTypeBytes() // todo fix this
	case typeBoolean:
		return builder.FieldTypeBool()
	}
	return builder.FieldTypeString() // todo, probably fatal instead
}

func main() {
	schemaFile, err := os.Open("schema.json")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	defer schemaFile.Close()

	schemaByteValue, err := ioutil.ReadAll(schemaFile)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	var schema configstoreSchema
	err = json.Unmarshal(schemaByteValue, &schema)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	var messages []*builder.MessageBuilder
	var services []*builder.ServiceBuilder

	for name, kind := range schema.Kinds {
		// Build the message descriptor for the entity itself
		message := builder.NewMessage(name)
		for _, field := range kind.Fields {
			message.AddField(
				builder.NewField(
					field.Name,
					convertToType(field.Type),
				).
					SetNumber(field.ID).
					SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" %s", field.Comment)}),
			)
		}
		messages = append(messages, message)

		// Build the request-response messages for the List method
		listRequestMessage := builder.NewMessage(fmt.Sprintf("List%sRequest", name)).
			AddField(builder.NewField("start", builder.FieldTypeBytes()).SetComments(builder.Comments{LeadingComment: " The start cursor from a previous List call, or null"})).
			AddField(builder.NewField("limit", builder.FieldTypeUInt32()).SetComments(builder.Comments{LeadingComment: " The maximum number of results to return, or null for no limit"}))
		listResponseMessage := builder.NewMessage(fmt.Sprintf("List%sResponse", name)).
			AddField(builder.NewField("next", builder.FieldTypeBytes()).SetComments(builder.Comments{LeadingComment: " The cursor to pass to the start field of the next List call"})).
			AddField(builder.NewField("moreResults", builder.FieldTypeBool()).SetRepeated().SetComments(builder.Comments{LeadingComment: " True if there are more results available in a future List call"})).
			AddField(builder.NewField("entities", builder.FieldTypeMessage(message)).SetRepeated().SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The paginated list of %ss", name)}))

		// Build the request-response message for the Get method
		getRequestMessage := builder.NewMessage(fmt.Sprintf("Get%sRequest", name)).
			AddField(builder.NewField("id", builder.FieldTypeUInt64()).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The ID of the %s to load", name)}))
		getResponseMessage := builder.NewMessage(fmt.Sprintf("Get%sResponse", name)).
			AddField(builder.NewField("entity", builder.FieldTypeMessage(message)).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The %s that was fetched, or null if it didn't exist", name)}))

		// Build the request-response message for the Update method
		updateRequestMessage := builder.NewMessage(fmt.Sprintf("Update%sRequest", name)).
			AddField(builder.NewField("entity", builder.FieldTypeMessage(message)).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The %s entity to update", name)}))
		updateResponseMessage := builder.NewMessage(fmt.Sprintf("Update%sResponse", name)).
			AddField(builder.NewField("entity", builder.FieldTypeMessage(message)).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The stored version of the %s entity", name)}))

		// Build the request-response message for the Create method
		createRequestMessage := builder.NewMessage(fmt.Sprintf("Create%sRequest", name)).
			AddField(builder.NewField("entity", builder.FieldTypeMessage(message)).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The %s entity to create; if %s uses auto-generated IDs, the ID field is ignored", name, name)}))
		createResponseMessage := builder.NewMessage(fmt.Sprintf("Create%sResponse", name)).
			AddField(builder.NewField("entity", builder.FieldTypeMessage(message)).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The stored version of the %s entity", name)}))

		// Build the request-response message for the Delete method
		deleteRequestMessage := builder.NewMessage(fmt.Sprintf("Delete%sRequest", name)).
			AddField(builder.NewField("id", builder.FieldTypeUInt64()).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The ID of the %s to delete", name)}))
		deleteResponseMessage := builder.NewMessage(fmt.Sprintf("Delete%sResponse", name)).
			AddField(builder.NewField("entity", builder.FieldTypeMessage(message)).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The version of the %s entity that was deleted", name)}))

		messages = append(messages, listRequestMessage)
		messages = append(messages, listResponseMessage)
		messages = append(messages, getRequestMessage)
		messages = append(messages, getResponseMessage)
		messages = append(messages, updateRequestMessage)
		messages = append(messages, updateResponseMessage)
		messages = append(messages, createRequestMessage)
		messages = append(messages, createResponseMessage)
		messages = append(messages, deleteRequestMessage)
		messages = append(messages, deleteResponseMessage)

		service := builder.NewService(fmt.Sprintf("%sService", name)).
			AddMethod(builder.NewMethod(
				"List",
				builder.RpcTypeMessage(listRequestMessage, false),
				builder.RpcTypeMessage(listResponseMessage, false),
			).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" Fetch a page of %s entities", name)})).
			AddMethod(builder.NewMethod(
				"Get",
				builder.RpcTypeMessage(getRequestMessage, false),
				builder.RpcTypeMessage(getResponseMessage, false),
			).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" Retrieve a single %s, if it exists", name)})).
			AddMethod(builder.NewMethod(
				"Update",
				builder.RpcTypeMessage(updateRequestMessage, false),
				builder.RpcTypeMessage(updateResponseMessage, false),
			).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" Update a single %s", name)})).
			AddMethod(builder.NewMethod(
				"Create",
				builder.RpcTypeMessage(createRequestMessage, false),
				builder.RpcTypeMessage(createResponseMessage, false),
			).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" Create a single %s", name)})).
			AddMethod(builder.NewMethod(
				"Delete",
				builder.RpcTypeMessage(deleteRequestMessage, false),
				builder.RpcTypeMessage(deleteResponseMessage, false),
			).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" Delete a single %s", name)}))
		services = append(services, service)
	}

	fileBuilder := builder.NewFile("")
	fileBuilder.SetProto3(true)
	fileBuilder.SetPackageName(schema.Name)
	for _, message := range messages {
		fileBuilder.AddMessage(message)
	}
	for _, service := range services {
		fileBuilder.AddService(service)
	}
	fileDesc, err := fileBuilder.Build()

	printer := new(protoprint.Printer)
	protoFile, err := os.Create("testclient/testclient.proto")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	defer protoFile.Close()

	fmt.Println("writing testclient/testclient.proto...")
	err = printer.PrintProtoFile(fileDesc, protoFile)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	lis, err := net.Listen("tcp", ":13389")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	grpcServer := grpc.NewServer()
	emptyServer := new(emptyServerInterface)
	for _, service := range services {
		grpcServer.RegisterService(
			&grpc.ServiceDesc{
				ServiceName: fmt.Sprintf("%s.%s", schema.Name, service.GetName()),
				HandlerType: (*emptyServerInterface)(nil),
				Methods: []grpc.MethodDesc{
					{
						MethodName: "List",
						Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
							return nil, fmt.Errorf("not implemented")
						},
					},
					{
						MethodName: "Get",
						Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
							return nil, fmt.Errorf("not implemented")
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
