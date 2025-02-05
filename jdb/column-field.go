package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type Field struct {
	Owner      interface{}
	Column     *Column
	Schema     string
	Table      string
	As         string
	Field      string
	Name       string
	Atrib      string
	Agregation TypeAgregation
	Alias      string
	Value      interface{}
}

func (s *Field) Define() et.Json {
	return et.Json{
		"schema": s.Schema,
		"table":  s.Table,
		"as":     s.As,
		"field":  s.Field,
		"name":   s.Name,
		"atrib":  s.Atrib,
		"alias":  s.Alias,
		"value":  s.Value,
	}
}

/**
* TableName
* @return string
**/
func (s *Field) TableName() string {
	return strs.Format("%s.%s", s.Table, s.Name)
}

/**
* TableField
* @return string
**/
func (s *Field) TableField() string {
	result := ""
	result = strs.Append(result, s.Schema, "")
	result = strs.Append(result, s.Table, ".")
	result = strs.Append(result, s.Field, ".")
	result = strs.Append(result, s.Atrib, ".")

	return result
}

/**
* AsField
* @return string
**/
func (s *Field) AsField() string {
	result := ""
	result = strs.Append(result, s.As, "")
	if s.Field != s.Name {
		result = strs.Append(result, s.Field, ".")
	}
	result = strs.Append(result, s.Name, ".")

	return result
}

/**
* Caption
* @return string
**/
func (s *Field) Caption() string {
	return s.Name
}
