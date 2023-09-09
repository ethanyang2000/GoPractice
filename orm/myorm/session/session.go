package session

import (
	"database/sql"
	"orm/myorm/clause"
	"orm/myorm/dialect"
	"orm/myorm/log"
	"orm/myorm/schema"
	"strings"
)

type Session struct {
	db       *sql.DB
	sql      strings.Builder
	sqlValue []interface{}
	dialect  dialect.Dialect
	refTable *schema.Schema
	clause   *clause.Clause
	tx       *sql.Tx
}

type DataBase interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
}

func NewSession(d *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      d,
		dialect: dialect,
		clause:  clause.NewClause(),
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlValue = s.sqlValue[0:0]
}

func (s *Session) DB() DataBase {
	if s.tx == nil {
		return s.db
	}
	return s.tx
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
