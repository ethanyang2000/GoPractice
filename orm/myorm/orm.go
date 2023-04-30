package myorm

import (
	"database/sql"
	"orm/myorm/log"
	"orm/myorm/session"
)

type Engine struct {
	db *sql.DB
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

	eng := &Engine{
		db: db,
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
	return session.NewSession(e.db)
}
