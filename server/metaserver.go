package main

import "context"

type configstoreMetaServiceServer struct {
	schema *configstoreSchema
}

func convertType(t configstoreSchemaKindFieldType) ValueType {
	switch t {
	case typeDouble:
		return ValueType_Double
	case typeInt64:
		return ValueType_Int64
	case typeString:
		return ValueType_String
	case typeTimestamp:
		return ValueType_Timestamp
	case typeBoolean:
		return ValueType_Boolean
	}

	return ValueType_Double
}

func (s *configstoreMetaServiceServer) GetSchema(ctx context.Context, req *GetSchemaRequest) (*GetSchemaResponse, error) {
	kinds := make([]*Kind, 0)
	for kindName, kind := range s.schema.Kinds {
		fields := make([]*Field, 0)
		for _, field := range kind.Fields {
			fields = append(fields, &Field{
				Id:      field.ID,
				Name:    field.Name,
				Type:    convertType(field.Type),
				Comment: field.Comment,
			})
		}

		kinds = append(kinds, &Kind{
			Name:   kindName,
			Fields: fields,
		})
	}

	schema := &Schema{
		Name:  s.schema.Name,
		Kinds: kinds,
	}

	return &GetSchemaResponse{
		Schema: schema,
	}, nil
}
