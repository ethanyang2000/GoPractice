package session

import (
	"database/sql"
	"fmt"
	"orm/myorm/clause"
	"orm/myorm/log"
	"orm/myorm/schema"
	"reflect"
	"strings"
)

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
	// value is a ptr to the object
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.RefTable().Model) {
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
	for _, name := range s.refTable.FieldNames {
		field, _ := s.RefTable().GetField(name)
		s := fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag)
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

func (s *Session) Insert(vars ...interface{}) (int64, error) {
	// var is a ptr to the object
	colValues := make([]interface{}, 0)
	s.Model(vars[0])
	s.clause.Set(clause.INSERT, s.RefTable().Name, s.RefTable().FieldNames)
	for _, obj := range vars {
		if reflect.TypeOf(obj) != reflect.TypeOf(s.RefTable().Model) {
			log.Error("unmatched objects for INSERT")
			return 0, fmt.Errorf("unmatched objects")
		}
		s.callHook(BeforeInsert, obj)
		colValues = append(colValues, s.RefTable().Flatten(obj))
	}
	s.clause.Set(clause.VALUES, colValues...)
	sql, sqlVars := s.clause.Build(clause.INSERT, clause.VALUES)
	r, err := s.Raw(sql, sqlVars...).Exec()
	if err != nil {
		log.Error(err)
	}
	return r.RowsAffected()
}

func (s *Session) Find(target interface{}) error {
	// target is a ptr to a slice of objects
	sliceValue := reflect.Indirect(reflect.ValueOf(target))
	objType := sliceValue.Type().Elem()
	newInstanceValue := reflect.Indirect(reflect.New(objType))
	s.Model(newInstanceValue.Addr().Interface())
	s.callHook(BeforeQuery, newInstanceValue.Addr().Interface())
	s.clause.Set(clause.SELECT, s.RefTable().Name, s.RefTable().FieldNames)
	sql, sqlVars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	results, err := s.Raw(sql, sqlVars...).QueryRows()
	if err != nil {
		log.Error(err)
		return err
	}
	defer results.Close()

	for results.Next() {
		newObjValue := reflect.Indirect(reflect.New(objType))
		addr := []interface{}{}
		for _, name := range s.RefTable().FieldNames {
			addr = append(addr, newObjValue.FieldByName(name).Addr().Interface())
		}
		err := results.Scan(addr...)
		if err != nil {
			log.Error(err)
			return err
		}
		sliceValue.Set(reflect.Append(sliceValue, newObjValue))
		s.callHook(AfterQuery, newObjValue.Addr().Interface())
	}
	return nil
}

func (s *Session) Update(kv ...interface{}) (int64, error) {
	kvMap, ok := kv[1].(map[string]interface{})
	if !ok {
		if len(kv)%2 != 0 {
			log.Error("illegal input for update")
			return 0, fmt.Errorf("illegal input for update")
		}
		kvMap = make(map[string]interface{})
		for idx := 0; idx < len(kv); idx += 2 {
			kvMap[kv[idx].(string)] = kv[idx+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, kvMap)
	sql, sqlVars := s.clause.Build(clause.UPDATE, clause.WHERE)
	r, err := s.Raw(sql, sqlVars...).Exec()
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return r.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, sqlVars := s.clause.Build(clause.DELETE, clause.WHERE)
	r, err := s.Raw(sql, sqlVars...).Exec()
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return r.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, sqlVars := s.clause.Build(clause.COUNT, clause.WHERE)
	r := s.Raw(sql, sqlVars...).QueryRow()
	var res int64
	if err := r.Scan(&res); err != nil {
		log.Error(err)
		return 0, err
	}
	return res, nil
}

func (s *Session) Where(vars ...interface{}) *Session {
	// [conditions, var1, var2...]
	s.clause.Set(clause.WHERE, vars...)
	return s
}

func (s *Session) Limit(inp int) *Session {
	s.clause.Set(clause.LIMIT, inp)
	return s
}

func (s *Session) OrderBy(vars ...interface{}) *Session {
	// [orderBy1, orderBy2...]
	s.clause.Set(clause.ORDERBY, vars...)
	return s
}

func (s *Session) First(target interface{}) error {
	// target is the ptr to the object
	targetValue := reflect.Indirect(reflect.ValueOf(target))
	s.Model(target)
	s.callHook(BeforeQuery, target)
	s.clause.Set(clause.SELECT, s.RefTable().Name, s.RefTable().FieldNames)
	s.clause.Set(clause.LIMIT, 1)
	sql, sqlVars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	r := s.Raw(sql, sqlVars...).QueryRow()
	attr := []interface{}{}
	for _, name := range s.RefTable().FieldNames {
		attr = append(attr, targetValue.FieldByName(name).Addr().Interface())
	}
	if err := r.Scan(attr...); err != nil {
		log.Error(err)
		return err
	}
	s.callHook(AfterQuery, target)
	return nil
}
