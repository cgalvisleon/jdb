package jdb

import "github.com/cgalvisleon/et/et"

type Dictionary struct {
	Column  *Column
	Name    string
	Key     string
	Value   interface{}
	Columns []*Column
}

/**
* NewDictionary
* @param key string
* @param value interface{}
* @return *Dictionary
**/
func NewDictionary(model *Model, name, key string, value interface{}) *Dictionary {
	col := model.GetColumn(key)
	return &Dictionary{
		Column:  col,
		Name:    name,
		Key:     key,
		Value:   value,
		Columns: []*Column{},
	}
}

func (s *Dictionary) Describe() et.Json {
	result := map[string]interface{}{
		"name":  s.Name,
		"key":   s.Key,
		"value": s.Value,
	}

	columns := []map[string]interface{}{}
	for _, col := range s.Columns {
		columns = append(columns, col.Describe())
	}
	result["columns"] = columns

	return result
}

/**
* DefineAtribute
* @param name string
* @param typeData TypeData
* @param def interface{}
* @return *Model
**/
func (s *Dictionary) DefineAtribute(name string, typeData TypeData, def interface{}) *Column {
	col := newColumn(s.Column.Model, name, "", TpAtribute, typeData, def)
	s.Columns = append(s.Columns, col)

	return col
}
