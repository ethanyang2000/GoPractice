package schema

import (
	"orm/myorm/dialect"
	"orm/myorm/log"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	// model is the ptr to an empty object
	Model interface{}
	Name  string
	// as map in go is unordered, we need a array to keep the fields in order
	FieldNames []string
	FieldMap   map[string]*Field
}

func (s *Schema) GetField(name string) (*Field, bool) {
	v, ok := s.FieldMap[name]
	return v, ok
}

func Parse(data interface{}, d dialect.Dialect) *Schema {
	// data is a ptr to the object
	log.Info("Parse called")
	modelType := reflect.Indirect(reflect.ValueOf(data)).Type()
	schema := &Schema{
		Model:    reflect.New(modelType).Interface(),
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
		schema.FieldNames = append(schema.FieldNames, f.Name)
	}
	return schema
}

func (s *Schema) Flatten(obj interface{}) []interface{} {
	// obj is the ptr to the object
	objValue := reflect.Indirect(reflect.ValueOf(obj))
	objVars := []interface{}{}
	for _, name := range s.FieldNames {
		field, _ := s.GetField(name)
		objVars = append(objVars, objValue.FieldByName(field.Name).Interface())
	}
	return objVars
}
