package session

import (
	"orm/myorm/log"
	"reflect"
)

type hookType int

const (
	AfterQuery hookType = iota
	BeforeQuery
	BeforeInsert
)

type AfterQueryHook interface {
	AfterQuery(s *Session) error
}

type BeforeQueryHook interface {
	BeforeQuery(s *Session) error
}

type BeforeInsertHook interface {
	BeforeInsert(s *Session) error
}

func (s *Session) callHook(t hookType, value interface{}) {
	// value should be the address of target object
	objType := reflect.ValueOf(value).Type()
	switch t {
	case BeforeQuery:
		iType := reflect.TypeOf((*BeforeQueryHook)(nil)).Elem()
		if objType.Implements(iType) {
			err := value.(BeforeQueryHook).BeforeQuery(s)
			if err != nil {
				log.Error(err)
			}
		}
	case AfterQuery:
		iType := reflect.TypeOf((*AfterQueryHook)(nil)).Elem()
		if objType.Implements(iType) {
			err := value.(AfterQueryHook).AfterQuery(s)
			if err != nil {
				log.Error(err)
			}
		}
	case BeforeInsert:
		iType := reflect.TypeOf((*BeforeInsertHook)(nil)).Elem()
		if objType.Implements(iType) {
			err := value.(BeforeInsertHook).BeforeInsert(s)
			if err != nil {
				log.Error(err)
			}
		}
	default:
		log.Error("invalid hook type")
	}
}
