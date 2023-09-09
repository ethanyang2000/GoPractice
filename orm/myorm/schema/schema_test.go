package schema

import (
	"orm/myorm/dialect"
	"testing"
)

func TestParse(t *testing.T) {
	type User struct {
		Name  string
		Score float32
		Age   int `myorm:"PRIMARY KEY"`
	}
	user := &User{
		Name:  "Tom",
		Score: 67.5,
		Age:   12,
	}
	dia, _ := dialect.GetDialect("sqlite3")
	schema := Parse(user, dia)
	if schema.Name != "User" || len(schema.FieldMap) != 3 {
		t.Fatal("failed to parse User struct")
	}
	if nameField, ok := schema.GetField("Name"); !ok || nameField.Type != "text" || nameField.Name != "Name" {
		t.Fatal("failed to parse User struct")
	}
	v, ok := schema.GetField("Age")
	if !ok || v.Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
	vars := schema.Flatten(user)
	if len(vars) != 3 || vars[0].(string) != "Tom" || vars[1].(float32) != 67.5 || vars[2].(int) != 12 {
		t.Fatal("failed to flatten object")
	}
}
