package main

import (
	"fmt"
	"orm/myorm"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func main() {
	engine, _ := myorm.NewEngine("sqlite3", "gee.db")
	defer engine.Close()
	s := engine.NewSession()
	s.Model(&User{})
	s.DropTable()
	s.CreateTable()
	if !s.HasTable() {
		fmt.Printf("Failed to create table User")
	}
}
