package session

import (
	"database/sql"
	"fmt"
	"orm/myorm/dialect"
	"orm/myorm/log"
	"orm/myorm/schema"
	"reflect"
	"strings"
)

type Session struct {
	db       *sql.DB
	sql      strings.Builder
	sqlValue []interface{}
	dialect  dialect.Dialect
	refTable *schema.Schema
}

func NewSession(d *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      d,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlValue = s.sqlValue[0:0]
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlValue = append(s.sqlValue, values...)
	return s
}

func (s *Session) Exec() (sql.Result, error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlValue)
	res, err := s.DB().Exec(s.sql.String(), s.sqlValue...)
	if err != nil {
		log.Error(err)
	}
	return res, err
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlValue)
	return s.DB().QueryRow(s.sql.String(), s.sqlValue...)
}

func (s *Session) QueryRows() (*sql.Rows, error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlValue)
	rows, err := s.DB().Query(s.sql.String(), s.sqlValue...)
	if err != nil {
		log.Error(err)
	}
	return rows, err
}

func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("model has not been set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.refTable
	column := make([]string, 0)
	for _, v := range s.refTable.FieldMap {
		s := fmt.Sprintf("%s %s %s", v.Name, v.Type, v.Tag)
		column = append(column, s)
	}
	str := strings.Join(column, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, str)).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
