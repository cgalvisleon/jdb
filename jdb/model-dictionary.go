package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/et"
)

type Dictionary struct {
	Column       *Column
	Key          string
	Value        interface{}
	Columns      []*Column
	Dictionaries map[string][]*Dictionary
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
		Column:       col,
		Key:          key,
		Value:        value,
		Columns:      []*Column{},
		Dictionaries: map[string][]*Dictionary{},
	}
}

/**
* Describe
* @return et.Json
**/
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
* DefineDictionary
* @param name string
* @param key string
* @param value interface{}
* @return *Dictionary
**/
func (s *Dictionary) DefineDictionary(key, value string) *Dictionary {
	results := s.Dictionaries[value]
	idx := slices.IndexFunc(results, func(e *Dictionary) bool { return e.Key == key })
	if idx != -1 {
		return results[idx]
	}

	result := NewDictionary(s.Column.Model, key, value)
	s.Dictionaries[value] = append(results, result)

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

/**
* DefineList
* @param name string
* @param typeData TypeData
* @param def interface{}
* @return *Model
**/
func (s *Dictionary) DefineList(name string, typeData TypeData, values []interface{}, def interface{}) *Column {
	col := newColumn(s.Column.Model, name, "", TpAtribute, typeData, def)
	col.Values = values
	s.Columns = append(s.Columns, col)

	return col
}
