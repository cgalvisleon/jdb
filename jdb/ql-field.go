package jdb

import (
	"regexp"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type TypeAgregation int

const (
	Nag TypeAgregation = iota
	AgregationSum
	AgregationCount
	AgregationAvg
	AgregationMin
	AgregationMax
)

func (s TypeAgregation) Str() string {
	switch s {
	case AgregationSum:
		return "SUM"
	case AgregationCount:
		return "COUNT"
	case AgregationAvg:
		return "AVG"
	case AgregationMin:
		return "MIN"
	case AgregationMax:
		return "MAX"
	}

	return ""
}

type Agregation struct {
	Agregation string
	pattern    string
	re         *regexp.Regexp
}

var agregations = map[TypeAgregation]*Agregation{
	Nag:             {Agregation: "", pattern: ""},
	AgregationSum:   {Agregation: "SUM", pattern: `SUM\([a-zA-Z0-9_]+\)$`},
	AgregationCount: {Agregation: "COUNT", pattern: `COUNT\([a-zA-Z0-9_]+\)$`},
	AgregationAvg:   {Agregation: "AVG", pattern: `AVG\([a-zA-Z0-9_]+\)$`},
	AgregationMin:   {Agregation: "MIN", pattern: `MIN\([a-zA-Z0-9_]+\)$`},
	AgregationMax:   {Agregation: "MAX", pattern: `MAX\([a-zA-Z0-9_]+\)$`},
}

type TypeResult int

const (
	TpResult TypeResult = iota
	TpList
)

/**
* StrToTypeResult
* @param str string
* @return TypeResult
**/
func StrToTypeResult(str string) TypeResult {
	switch str {
	case "list":
		return TpList
	}

	return TpResult
}

/**
* init
**/
func init() {
	for _, agregation := range agregations {
		re, err := regexp.Compile(agregation.pattern)
		if err != nil {
			continue
		}
		agregation.re = re
	}
}

type Field struct {
	Column     *Column
	Schema     string
	Table      string
	As         string
	Name       string
	Source     string
	Agregation TypeAgregation
	Value      interface{}
	Alias      string
	Hidden     bool
	Page       int
	Rows       int
	TpResult   TypeResult
}

/**
* Describe
* @return et.Json
**/
func (s *Field) Describe() et.Json {
	result, err := et.Object(s)
	if err != nil {
		return et.Json{}
	}

	return result
}

/**
* NewField
* @param column *Column
* @return *Field
**/
func NewField(column *Column) *Field {
	result := &Field{}
	if column == nil {
		return result
	}

	result.Column = column
	result.Name = column.Name
	result.Alias = column.Name
	result.Hidden = column.Hidden

	if column.TypeColumn == TpRelatedTo {
		result.Page = 1
		result.Rows = 30
		result.TpResult = TpResult
	}

	if column.Model == nil {
		return result
	}

	result.Schema = column.Model.Schema.Name
	result.Table = column.Model.Name

	if column.Source == nil {
		return result
	}

	result.Source = column.Source.Name

	return result
}

func (s *Field) Define() et.Json {
	return et.Json{
		"schema":     s.Schema,
		"table":      s.Table,
		"as":         s.As,
		"name":       s.Name,
		"source":     s.Source,
		"agregation": s.Agregation.Str(),
		"alias":      s.Alias,
		"value":      s.Value,
	}
}

func (s *Field) SetAgregation(agr TypeAgregation) {
	s.Agregation = agr
	switch agr {
	case AgregationSum:
		s.Alias = strs.Format("sum_%s", s.Name)
	case AgregationCount:
		s.Alias = strs.Format("count_%s", s.Name)
	case AgregationAvg:
		s.Alias = strs.Format("avg_%s", s.Name)
	case AgregationMin:
		s.Alias = strs.Format("min_%s", s.Name)
	case AgregationMax:
		s.Alias = strs.Format("max_%s", s.Name)
	}
}

/**
* AsTable
* @return string
**/
func (s *Field) AsTable() string {
	return strs.Format("%s.%s", s.Table, s.Name)
}

/**
* AsField
* @return string
**/
func (s *Field) AsField() string {
	result := ""
	result = strs.Append(result, s.Schema, "")
	result = strs.Append(result, s.Table, ".")
	result = strs.Append(result, s.Source, ".")
	result = strs.Append(result, s.Name, ".")

	return result
}

/**
* AsName
* @return string
**/
func (s *Field) AsName() string {
	if s.As != "" {
		return strs.Format("%s.%s", s.As, s.Name)
	}

	return s.Name
}

/**
* GetField
* @return *Field
**/
func (s *Column) GetField() *Field {
	result := NewField(s)

	return result
}
