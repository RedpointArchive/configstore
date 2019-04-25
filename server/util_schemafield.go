package main

func findSchemaFieldByID(schemaKind *SchemaKind, id int32) *SchemaField {
	for _, field := range schemaKind.Fields {
		if field.Id == id {
			return field
		}
	}

	return nil
}
