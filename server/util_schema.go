package main

import (
	"fmt"
)

func findSchemaKindByName(schema *Schema, name string) (*SchemaKind, error) {
	for kindName, kind := range schema.Kinds {
		if kindName == name {
			return kind, nil
		}
	}
	return nil, fmt.Errorf("no such kind '%s'", name)
}
