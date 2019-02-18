package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang/protobuf/jsonpb"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

func convertToType(
	t ValueType,
	keyMessage *builder.MessageBuilder,
	timestampMessage *desc.MessageDescriptor,
) *builder.FieldType {
	switch t {
	case ValueType_double:
		return builder.FieldTypeDouble()
	case ValueType_int64:
		return builder.FieldTypeInt64()
	case ValueType_uint64:
		return builder.FieldTypeUInt64()
	case ValueType_string:
		return builder.FieldTypeString()
	case ValueType_timestamp:
		return builder.FieldTypeImportedMessage(timestampMessage)
	case ValueType_boolean:
		return builder.FieldTypeBool()
	case ValueType_bytes:
		return builder.FieldTypeBytes()
	case ValueType_key:
		return builder.FieldTypeMessage(keyMessage)
	}
	fmt.Printf("fatal: no such field type '%s'", string(t))
	os.Exit(1)
	return builder.FieldTypeString() // never reached
}

type watchTypeEnumValues struct {
	Created *desc.EnumValueDescriptor
	Updated *desc.EnumValueDescriptor
	Deleted *desc.EnumValueDescriptor
}

type commonMessageDescriptors struct {
	Timestamp   *desc.MessageDescriptor
	PartitionId *desc.MessageDescriptor
	PathElement *desc.MessageDescriptor
	Key         *desc.MessageDescriptor
}

type generatorResult struct {
	Messages                 []*builder.MessageBuilder
	Services                 []*builder.ServiceBuilder
	FileBuilder              *builder.FileBuilder
	FileDesc                 *desc.FileDescriptor
	Schema                   *Schema
	KindMap                  map[*builder.ServiceBuilder]*SchemaKind
	KindNameMap              map[*builder.ServiceBuilder]string
	MessageMap               map[string]*desc.MessageDescriptor
	WatchTypeEnumValues      *watchTypeEnumValues
	CommonMessageDescriptors *commonMessageDescriptors
}

func generate(path string) (*generatorResult, error) {
	timestampFileDescriptor, err := desc.LoadFileDescriptor("google/protobuf/timestamp.proto")
	if err != nil {
		return nil, err
	}

	var timestampMessage *desc.MessageDescriptor
	for _, tfd := range timestampFileDescriptor.GetMessageTypes() {
		if tfd.GetName() == "Timestamp" {
			timestampMessage = tfd
			break
		}
	}
	if timestampMessage == nil {
		return nil, fmt.Errorf("unable to locate Timestamp proto descriptor in google/protobuf/timestamp.proto")
	}

	schemaFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer schemaFile.Close()

	var schema Schema
	err = jsonpb.Unmarshal(schemaFile, &schema)
	if err != nil {
		return nil, fmt.Errorf("unable to deserialize schema.json: %v", err)
	}

	var messages []*builder.MessageBuilder
	var services []*builder.ServiceBuilder
	var enums []*builder.EnumBuilder
	kindMap := make(map[*builder.ServiceBuilder]*SchemaKind)
	kindNameMap := make(map[*builder.ServiceBuilder]string)
	messageMap := make(map[string]*desc.MessageDescriptor)

	partitionIdMessage := builder.NewMessage("PartitionId")
	partitionIdMessage.AddField(
		builder.NewField("namespace", builder.FieldTypeString()).
			SetNumber(1),
	)
	messages = append(messages, partitionIdMessage)

	pathElementIdMessage := builder.NewOneOf("idType")
	pathElementIdMessage.AddChoice(
		builder.NewField("id", builder.FieldTypeInt64()).
			SetNumber(2).
			SetOptional().
			SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The ID in the path component")}),
	)
	pathElementIdMessage.AddChoice(
		builder.NewField("name", builder.FieldTypeString()).
			SetNumber(3).
			SetOptional().
			SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The name in the path component")}),
	)

	pathElementMessage := builder.NewMessage("PathElement")
	pathElementMessage.AddField(
		builder.NewField("kind", builder.FieldTypeString()).
			SetNumber(1).
			SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The kind in the path component")}),
	)
	pathElementMessage.AddOneOf(pathElementIdMessage)
	messages = append(messages, pathElementMessage)

	keyMessage := builder.NewMessage("Key")
	keyMessage.AddField(
		builder.NewField("partitionId", builder.FieldTypeMessage(partitionIdMessage)).
			SetNumber(1).
			SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The partition that the entity is stored in; if omitted, the default is used")}),
	)
	keyMessage.AddField(
		builder.NewField("path", builder.FieldTypeMessage(pathElementMessage)).
			SetNumber(2).
			SetRepeated().
			SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The path to the entity")}),
	)
	messages = append(messages, keyMessage)

	partitionIdDescriptor, err := partitionIdMessage.Build()
	if err != nil {
		return nil, err
	}
	pathElementDescriptor, err := pathElementMessage.Build()
	if err != nil {
		return nil, err
	}
	keyDescriptor, err := keyMessage.Build()
	if err != nil {
		return nil, err
	}

	common := &commonMessageDescriptors{
		Timestamp:   timestampMessage,
		PartitionId: partitionIdDescriptor,
		PathElement: pathElementDescriptor,
		Key:         keyDescriptor,
	}

	created := builder.NewEnumValue("Created").SetNumber(0)
	updated := builder.NewEnumValue("Updated").SetNumber(1)
	deleted := builder.NewEnumValue("Deleted").SetNumber(2)

	watchEventTypeEnum := builder.NewEnum("WatchEventType").
		AddValue(created).
		AddValue(updated).
		AddValue(deleted)
	enums = append(enums, watchEventTypeEnum)

	createdBuilt, err := created.Build()
	if err != nil {
		return nil, err
	}
	updatedBuilt, err := updated.Build()
	if err != nil {
		return nil, err
	}
	deletedBuilt, err := deleted.Build()
	if err != nil {
		return nil, err
	}

	watchTypeEnumValues := &watchTypeEnumValues{
		Created: createdBuilt,
		Updated: updatedBuilt,
		Deleted: deletedBuilt,
	}

	for name, kind := range schema.Kinds {
		// Build the message descriptor for the entity itself
		message := builder.NewMessage(name)
		message.AddField(
			builder.NewField("key", builder.FieldTypeMessage(keyMessage)).
				SetNumber(1).
				SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The key of the %s", name)}),
		)
		for _, field := range kind.Fields {
			if field.Id == 1 {
				log.Fatalln("unexpected ID 1 in kind field; IDs must start at 2")
			}
			message.AddField(
				builder.NewField(
					field.Name,
					convertToType(field.Type, keyMessage, timestampMessage),
				).
					SetNumber(field.Id).
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
			AddField(builder.NewField("moreResults", builder.FieldTypeBool()).SetComments(builder.Comments{LeadingComment: " True if there are more results available in a future List call"})).
			AddField(builder.NewField("entities", builder.FieldTypeMessage(message)).SetRepeated().SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The paginated list of %ss", name)}))

		// Build the request-response message for the Get method
		getRequestMessage := builder.NewMessage(fmt.Sprintf("Get%sRequest", name)).
			AddField(builder.NewField("key", builder.FieldTypeMessage(keyMessage)).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The ID of the %s to load", name)}))
		getResponseMessage := builder.NewMessage(fmt.Sprintf("Get%sResponse", name)).
			AddField(builder.NewField("entity", builder.FieldTypeMessage(message)).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The %s that was fetched, or null if it didn't exist", name)}))

		// Build the request-response message for the Watch method
		watchRequestMessage := builder.NewMessage(fmt.Sprintf("Watch%sRequest", name))
		watchEventMessage := builder.NewMessage(fmt.Sprintf("Watch%sEvent", name)).
			AddField(builder.NewField("type", builder.FieldTypeEnum(watchEventTypeEnum)).SetComments(builder.Comments{LeadingComment: " The type of modification"})).
			AddField(builder.NewField("entity", builder.FieldTypeMessage(message)).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The %s that was created, modified or deleted", name)})).
			AddField(builder.NewField("oldIndex", builder.FieldTypeInt64()).SetComments(builder.Comments{LeadingComment: " The old index of the entity in the collection, or -1 if it wasn't present"})).
			AddField(builder.NewField("newIndex", builder.FieldTypeInt64()).SetComments(builder.Comments{LeadingComment: " The new index of the entity in the collection, or -1 if it is no longer present"}))

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
			AddField(builder.NewField("key", builder.FieldTypeMessage(keyMessage)).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The ID of the %s to delete", name)}))
		deleteResponseMessage := builder.NewMessage(fmt.Sprintf("Delete%sResponse", name)).
			AddField(builder.NewField("entity", builder.FieldTypeMessage(message)).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" The version of the %s entity that was deleted", name)}))

		messages = append(messages, listRequestMessage)
		messages = append(messages, listResponseMessage)
		messages = append(messages, getRequestMessage)
		messages = append(messages, getResponseMessage)
		messages = append(messages, watchRequestMessage)
		messages = append(messages, watchEventMessage)
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
				"Watch",
				builder.RpcTypeMessage(watchRequestMessage, false),
				builder.RpcTypeMessage(watchEventMessage, true),
			).SetComments(builder.Comments{LeadingComment: fmt.Sprintf(" Watch all %s entities for changes", name)})).
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

		kindMap[service] = kind
		kindNameMap[service] = name
	}

	for _, message := range messages {
		msgDesc, err := message.Build()
		if err != nil {
			return nil, err
		}
		messageMap[message.GetName()] = msgDesc
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
	for _, enum := range enums {
		fileBuilder.AddEnum(enum)
	}
	fileDesc, err := fileBuilder.Build()
	if err != nil {
		return nil, err
	}

	return &generatorResult{
		Messages:                 messages,
		Services:                 services,
		FileBuilder:              fileBuilder,
		FileDesc:                 fileDesc,
		Schema:                   &schema,
		KindMap:                  kindMap,
		KindNameMap:              kindNameMap,
		MessageMap:               messageMap,
		WatchTypeEnumValues:      watchTypeEnumValues,
		CommonMessageDescriptors: common,
	}, nil
}
