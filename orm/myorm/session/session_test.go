package session

import (
	"database/sql"
	"orm/myorm/dialect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id   int
	Name string
	Age  int
}

var (
	user1 = User{1, "Tom", 18}
	user2 = User{2, "Tom", 25}
	user3 = User{3, "Jack", 25}
)

func TestSession(t *testing.T) {
	u := User{}
	db, err := sql.Open("sqlite3", "gee.db")
	if err != nil {
		t.Fatal(err)
	}
	s := NewSession(db, dialect.DialectMap["sqlite3"])
	s.Model(&u)
	if s.RefTable() == nil || (len(s.RefTable().FieldNames) != 3) {
		t.Fatal("model object failed")
	}
	s.DropTable()
	if s.CreateTable() != nil {
		t.Fatal(err)
	}
	result := ""
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='User'").Scan(&result)
	if err != nil {
		t.Fatal(err)
	}
	if result != "User" {
		t.Fatal("failed to create table")
	}
	if !s.HasTable() {
		t.Fatal("HasTable Query failed")
	}
	rows, insErr := s.Insert(&User{
		Id:   1,
		Name: "Tom",
		Age:  23,
	})
	if insErr != nil {
		t.Fatal(insErr)
	}
	if rows != 1 {
		t.Fatal("failed to insert one element")
	}
	QuerySlice := []User{}
	findErr := s.Find(&QuerySlice)
	if findErr != nil {
		t.Fatal("failed to find")
	}
	if len(QuerySlice) != 1 {
		t.Fatal(QuerySlice)
	}
	if QuerySlice[0].Id != 1 || QuerySlice[0].Name != "Tom" || QuerySlice[0].Age != 23 {
		t.Fatal("failed to find")
	}

	s.DropTable()
	s.CreateTable()
	s.Insert(&user1, &user2, &user3)

	res, e := s.Where("Name = ?", "Tom").Count()
	if e != nil || res != 2 {
		t.Fatal("failed to count")
	}

	var users []User
	err = s.Limit(1).Find(&users)
	if err != nil || len(users) != 1 {
		t.Fatal("failed to query with limit condition")
	}
	affected, _ := s.Where("Name = ? and Id = ?", "Tom", 1).Update("Age", 30)
	if affected != 1 {
		t.Fatal("failed to update row")
	}

	u_ptr := &User{}
	err = s.OrderBy("Age DESC").First(u_ptr)
	if err != nil {
		t.Fatal(err)
	}
	if u_ptr.Id != 1 || u_ptr.Age != 30 || u_ptr.Name != "Tom" {
		t.Fatal("failed to update")
	}

	err = s.OrderBy("Age ASC", "Id DESC").First(u_ptr)
	if err != nil {
		t.Fatal(err)
	}
	if u_ptr.Id != 3 || u_ptr.Age != 25 || u_ptr.Name != "Jack" {
		t.Fatal("failed to update")
	}

	affected, _ = s.Where("Name = ?", "Tom").Delete()
	count, _ := s.Count()
	if affected != 2 || count != 1 {
		t.Fatal("failed to delete or count")
	}
}
