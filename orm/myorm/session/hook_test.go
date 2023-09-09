package session

import (
	"database/sql"
	"orm/myorm/dialect"
	"orm/myorm/log"
	"reflect"
	"testing"
)

type Account struct {
	ID       int `geeorm:"PRIMARY KEY"`
	Password string
}

func (account *Account) BeforeInsert(s *Session) error {
	log.Info("before insert", account)
	account.ID += 1000
	return nil
}

func (account *Account) AfterQuery(s *Session) error {
	log.Info("after query", account)
	account.Password = "******"
	return nil
}

func TestSession_CallMethod(t *testing.T) {
	db, err := sql.Open("sqlite3", "gee.db")
	if err != nil {
		t.Fatal(err)
	}
	s := NewSession(db, dialect.DialectMap["sqlite3"])
	s.Model(&Account{})
	_ = s.DropTable()
	_ = s.CreateTable()
	a := Account{1, "123456"}
	log.Info(reflect.TypeOf(&a).NumMethod())
	_, _ = s.Insert(&Account{1, "123456"}, &Account{2, "qwerty"})

	u := &Account{}
	log.Info(reflect.TypeOf(u).NumMethod())
	err = s.First(u)
	if err != nil || u.ID != 1001 || u.Password != "******" {
		t.Fatal("Failed to call hooks after query, got", u)
	}
}
