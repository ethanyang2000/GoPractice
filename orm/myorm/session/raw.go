package session

import (
	"database/sql"
	"orm/myorm/log"
	"strings"
)

type Session struct {
	db       *sql.DB
	sql      strings.Builder
	sqlValue []interface{}
}

func NewSession(d *sql.DB) *Session {
	return &Session{
		db: d,
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
