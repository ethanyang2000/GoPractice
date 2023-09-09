package session

import "orm/myorm/log"

type transFunc func(*Session) (interface{}, error)

func (s *Session) Begin() (err error) {
	s.tx, err = s.db.Begin()
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *Session) Rollback() error {
	defer func() {
		s.tx = nil
	}()
	err := s.tx.Rollback()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s *Session) Commit() error {
	defer func() {
		s.tx = nil
	}()
	err := s.tx.Commit()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s *Session) Transection(fn transFunc) (res interface{}, err error) {
	err = s.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			s.Rollback()
			panic(p)
		} else if err != nil {
			log.Error(err)
			s.Rollback()
		} else {
			err = s.Commit()
			if err != nil {
				log.Error(err)
			}
		}
	}()
	res, err = fn(s)
	return
}
