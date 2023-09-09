package clause

import (
	"reflect"
	"testing"
)

func TestClause(t *testing.T) {
	c := NewClause()
	c.Set(LIMIT, 3)
	c.Set(SELECT, "User", []string{"name", "id"})
	c.Set(WHERE, "name = ?", "Tom")
	c.Set(ORDERBY, "id ASC", "name DESC")
	sql, vars := c.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT name,id FROM User WHERE name = ? ORDER BY id ASC,name DESC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}

	c.Set(INSERT, "User", []string{"id", "name"})
	c.Set(VALUES, []interface{}{3, "Tom"}, []interface{}{4, "Bob"})
	sql, vars = c.Build(INSERT, VALUES)
	t.Log(sql, vars)
	if sql != "INSERT INTO User (id,name) VALUES (?,?),(?,?)" {
		t.Fatal("failed to build sql")
	}
	if !reflect.DeepEqual(vars, []interface{}{3, "Tom", 4, "Bob"}) {
		t.Fatal("failed to build SQLVars")
	}
}
