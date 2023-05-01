package dialect

import "reflect"

var DialectMap map[string]Dialect = make(map[string]Dialect)

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	DialectMap[name] = dialect
}

func GetDialect(name string) (Dialect, bool) {
	v, ok := DialectMap[name]
	if !ok {
		return nil, ok
	}
	return v, ok
}
