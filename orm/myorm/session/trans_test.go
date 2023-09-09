package session

import (
	"database/sql"
	"orm/myorm/dialect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func rollBackTest(s *Session) (interface{}, error) {
	s.Insert(&user1)
	_, err := s.Insert(user2)
	return nil, err
}

func commitTest(s *Session) (interface{}, error) {
	s.Insert(&user1)
	s.Insert(&user2)
	s.Insert(&user3)
	dst := &[]User{}
	err := s.Where("Name = ?", "Tom").Find(dst)
	return dst, err
}

func TestTrans(t *testing.T) {
	db, err := sql.Open("sqlite3", "gee.db")
	if err != nil {
		t.Fatal(err)
	}
	s := NewSession(db, dialect.DialectMap["sqlite3"])
	s.Model(&User{})
	s.DropTable()
	s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}

	_, err1 := s.Transection(rollBackTest)
	if err1 == nil {
		t.Fatal("rollback test failed")
	}
	if num, _ := s.Count(); num != 0 {
		t.Fatal("rollback test failed")
	}

	s.DropTable()
	s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}

	res, err2 := s.Transection(commitTest)
	if err2 != nil {
		t.Fatal("commit test failed")
	}
	if num, _ := s.Count(); num != 3 {
		t.Fatal("commit test failed")
	}
	if r, ok := res.(*[]User); !ok || len(*r) != 2 || (*r)[0].Name != "Tom" || (*r)[1].Name != "Tom" {
		t.Fatal("commit test failed")
	}
}
