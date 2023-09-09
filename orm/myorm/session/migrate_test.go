package session

import (
	"database/sql"
	"orm/myorm/dialect"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrate(t *testing.T) {
	db, err := sql.Open("sqlite3", "gee.db")
	if err != nil {
		t.Fatal(err)
	}
	s := NewSession(db, dialect.DialectMap["sqlite3"])

	s.Raw("DROP TABLE IF EXISTS User").Exec()
	s.Raw("CREATE TABLE User(Name text PRIMARY KEY, XXX integer);").Exec()
	s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()

	rows, _ := s.Raw("SELECT * FROM User").QueryRows()
	columns, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Name", "XXX"}) {
		t.Fatal("Failed to build table User")
	}

	s.Migrate(&User{})

	rows, _ = s.Raw("SELECT * FROM User").QueryRows()
	columns, _ = rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Id", "Name", "Age"}) {
		t.Fatal("Failed to migrate table User, got columns", columns)
	}
}
