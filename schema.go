package main

type configstoreSchema struct {
	Name  string                           `json:"name"`
	Kinds map[string]configstoreSchemaKind `json:"kinds"`
}

type configstoreSchemaKind struct {
	Ancestors []string                     `json:"ancestors"`
	Fields    []configstoreSchemaKindField `json:"fields"`
}

type configstoreSchemaKindFieldType string

const (
	typeDouble    configstoreSchemaKindFieldType = "double"
	typeInt64     configstoreSchemaKindFieldType = "int64"
	typeString    configstoreSchemaKindFieldType = "string"
	typeTimestamp configstoreSchemaKindFieldType = "timestamp"
	typeBoolean   configstoreSchemaKindFieldType = "bool"
)

type configstoreSchemaKindField struct {
	ID      int32                          `json:"id"`
	Name    string                         `json:"name"`
	Type    configstoreSchemaKindFieldType `json:"type"`
	Comment string                         `json:"comment"`
}
