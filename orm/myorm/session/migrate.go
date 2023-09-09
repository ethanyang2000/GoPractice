package session

import (
	"fmt"
	"orm/myorm/log"
	"strings"
)

func diff(new, old []string) ([]string, []string) {
	add := []string{}
	del := []string{}
	memNew := map[string]bool{}
	memOld := map[string]bool{}
	for _, n := range new {
		memNew[n] = true
	}
	for _, o := range old {
		memOld[o] = true
		if _, ok := memNew[o]; !ok {
			del = append(del, o)
		}
	}
	for _, n := range new {
		if _, ok := memOld[n]; !ok {
			add = append(add, n)
		}
	}
	return add, del
}

func (s *Session) Migrate(value interface{}) (err error) {
	// value is a ptr to an object int the new table
	_, err = s.Transection(func(s *Session) (res interface{}, err error) {
		s.Model(value)
		if !s.HasTable() {
			log.Errorf("table %s do not exist", s.RefTable().Name)
		}
		r, err := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", s.RefTable().Name)).QueryRows()
		if err != nil {
			log.Error(err)
			return
		}
		cols, _ := r.Columns()
		add, del := diff(s.RefTable().FieldNames, cols)

		for _, a := range add {
			t, _ := s.RefTable().GetField(a)
			_, err = s.Raw(fmt.Sprintf("ALTER TABLE \"%s\" ADD COLUMN \"%s\" \"%s\"", s.RefTable().Name, t.Name, t.Type)).Exec()
			if err != nil {
				return
			}
		}

		if len(del) == 0 {
			return
		}

		newCols := strings.Join(s.refTable.FieldNames, ",")
		sqlStr := fmt.Sprintf("CREATE TABLE new_table AS SELECT %s from %s", newCols, s.RefTable().Name)
		_, err = s.Raw(sqlStr).Exec()
		if err != nil {
			return
		}

		s.DropTable()
		_, err = s.Raw(fmt.Sprintf("ALTER TABLE new_table RENAME TO %s", s.RefTable().Name)).Exec()
		if err != nil {
			return
		}
		return
	})
	return err
}
