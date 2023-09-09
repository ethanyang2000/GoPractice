package clause

import (
	"fmt"
	"strings"
)

type generator func(vars ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[LIMIT] = _limit
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

func _insert(vars ...interface{}) (string, []interface{}) {
	// [table_name, [fields]]
	// INSERT INTO $tableName ($fields)
	fields := strings.Join(vars[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", vars[0], fields), []interface{}{}
}

func _values(vars ...interface{}) (string, []interface{}) {
	// [[values of row1], []...]
	// VALUES ($v1), ($v2), ...
	sqlvars := []interface{}{}
	sampleVar := vars[0].([]interface{})
	sqlbackbone := "("
	for i := 0; i < len(sampleVar); i++ {
		if i == len(sampleVar)-1 {
			sqlbackbone += "?"
		} else {
			sqlbackbone += "?,"
		}
	}
	sqlbackbone += ")"
	sqlstr := "VALUES "
	for idx := 0; idx < len(vars); idx++ {
		if idx != 0 {
			sqlstr += ","
		}
		sqlstr += sqlbackbone
		sqlvars = append(sqlvars, vars[idx].([]interface{})...)
	}
	return sqlstr, sqlvars
}

func _select(vars ...interface{}) (string, []interface{}) {
	// [table_name, [fields]]
	// SELECT $fields FROM $tableName
	fields := strings.Join(vars[1].([]string), ",")
	return fmt.Sprintf("SELECT %s FROM %s", fields, vars[0]), []interface{}{}
}

func _where(vars ...interface{}) (string, []interface{}) {
	// WHERE $desc
	return fmt.Sprintf("WHERE %s", vars[0].(string)), vars[1:]
}

func _orderBy(vars ...interface{}) (string, []interface{}) {
	sqlStr := []string{}
	for _, col := range vars {
		sqlStr = append(sqlStr, col.(string))
	}
	return fmt.Sprintf("ORDER BY %s", strings.Join(sqlStr, ",")), []interface{}{}
}

func _limit(vars ...interface{}) (string, []interface{}) {
	return "LIMIT ?", vars
}

func _update(vars ...interface{}) (string, []interface{}) {
	// [tableName, map]
	// UPDATE table SET col=value, col=value
	kvMap := vars[1].(map[string]interface{})
	sqlVars := []interface{}{}
	sqlStr := fmt.Sprintf("UPDATE %s SET ", vars[0].(string))
	idx := 0
	for k, v := range kvMap {
		if idx == 0 {
			sqlStr += k + "=?"
		} else {
			sqlStr += ("," + k + "=?")
		}
		sqlVars = append(sqlVars, v)
	}
	return sqlStr, sqlVars
}

func _delete(vars ...interface{}) (string, []interface{}) {
	// DELETE FROM table
	return fmt.Sprintf("DELETE FROM %s ", vars[0]), []interface{}{}
}

func _count(vars ...interface{}) (string, []interface{}) {
	return _select(vars[0], []string{"COUNT(*)"})
}
