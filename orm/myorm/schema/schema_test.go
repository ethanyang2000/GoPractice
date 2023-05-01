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
	v, ok := schema.GetField("Age")
	if !ok || v.Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}
