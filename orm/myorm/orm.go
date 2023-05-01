package myorm

import (
	"database/sql"
	"fmt"
	"orm/myorm/dialect"
	"orm/myorm/log"
	"orm/myorm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (*Engine, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Error(err)
		return nil, err
	}

	d, ok := dialect.GetDialect(driver)
	if !ok {
		log.Error(fmt.Sprintf("dialect %s do not exist", driver))
		return nil, fmt.Errorf("dialect %s do not exist", driver)
	}

	eng := &Engine{
		db:      db,
		dialect: d,
	}

	log.Info("database connected")

	return eng, nil
}

func (e *Engine) Close() {
	err := e.db.Close()
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("database closed scuessfully")
}

func (e *Engine) NewSession() *session.Session {
	return session.NewSession(e.db, e.dialect)
}
