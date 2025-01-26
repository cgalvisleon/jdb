package jdb

import "github.com/cgalvisleon/et/et"

type Dictionary struct {
	Column     *Column
	Key        string
	Value      interface{}
	Columns    []*Column
	Dictionary map[string][]*Dictionary
}

/**
* NewDictionary
* @param key string
* @param value interface{}
* @return *Dictionary
**/
func NewDictionary(model *Model, key string, value interface{}) *Dictionary {
	col := model.GetColumn(key)
	return &Dictionary{
		Column:  col,
		Key:     key,
		Value:   value,
		Columns: []*Column{},
	}
}

func (s *Dictionary) Describe() et.Json {
	result := map[string]interface{}{
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
