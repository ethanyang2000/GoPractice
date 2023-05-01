package schema

import (
	"orm/myorm/dialect"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model interface{}
	Name  string
	//Fields   []*Field
	FieldMap map[string]*Field
}

func (s *Schema) GetField(name string) (*Field, bool) {
	v, ok := s.FieldMap[name]
	return v, ok
}

func Parse(data interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(data)).Type()
	schema := &Schema{
		Model:    data,
		Name:     modelType.Name(),
		FieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		f := modelType.Field(i)
		field := &Field{
			Name: f.Name,
			Type: d.DataTypeOf(reflect.Indirect(reflect.New(f.Type))),
		}
		if v, ok := f.Tag.Lookup("myorm"); ok {
			field.Tag = v
		}
		schema.FieldMap[f.Name] = field
	}
	return schema
}
