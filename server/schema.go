package main

type configstoreSchema struct {
	Name  string                           `json:"name"`
	Kinds map[string]configstoreSchemaKind `json:"kinds"`
}

type configstoreSchemaKind struct {
	Ancestors []string                     `json:"ancestors"`
	Fields    []configstoreSchemaKindField `json:"fields"`
	Editor    configstoreSchemaKindEditor  `json:"editor"`
}

type configstoreSchemaKindEditor struct {
	Singular string `json:"singular"`
	Plural   string `json:"plural"`
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
	ID      int32                            `json:"id"`
	Name    string                           `json:"name"`
	Type    configstoreSchemaKindFieldType   `json:"type"`
	Comment string                           `json:"comment"`
	Editor  configstoreSchemaKindFieldEditor `json:"editor"`
}

type configstoreSchemaKindFieldEditorType string

const (
	editorTypeDefault  configstoreSchemaKindFieldEditorType = ""
	editorTypePassword configstoreSchemaKindFieldEditorType = "password"
	editorTypeLookup   configstoreSchemaKindFieldEditorType = "lookup"
)

type configstoreSchemaKindFieldEditor struct {
	DisplayName string                               `json:"displayName"`
	Type        configstoreSchemaKindFieldEditorType `json:"type"`
	Readonly    bool                                 `json:"readonly"`
	ForeignType string                               `json:"foreignType"`
}
