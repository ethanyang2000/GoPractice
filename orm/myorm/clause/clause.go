package clause

import "strings"

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

func NewClause() *Clause {
	return &Clause{
		sql:     make(map[Type]string),
		sqlVars: make(map[Type][]interface{}),
	}
}

func (c *Clause) Set(t Type, vars ...interface{}) {
	sqlStr, sqlVar := generators[t](vars...)
	c.sql[t] = sqlStr
	c.sqlVars[t] = sqlVar
}

func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	sqlStr := []string{}
	sqlVars := []interface{}{}
	for _, order := range orders {
		if _, ok := c.sql[order]; ok {
			sqlStr = append(sqlStr, c.sql[order])
			sqlVars = append(sqlVars, c.sqlVars[order]...)
		}
	}
	c.sql = make(map[Type]string)
	c.sqlVars = make(map[Type][]interface{})
	return strings.Join(sqlStr, " "), sqlVars
}
