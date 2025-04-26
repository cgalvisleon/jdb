package jdb

import (
	"encoding/json"
	"regexp"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
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
	Column     *Column        `json:"column"`
	Schema     string         `json:"schema"`
	Table      string         `json:"table"`
	As         string         `json:"as"`
	Name       string         `json:"name"`
	Source     string         `json:"source"`
	Agregation TypeAgregation `json:"agregation"`
	Value      interface{}    `json:"value"`
	Alias      string         `json:"alias"`
	Hidden     bool           `json:"hidden"`
	Page       int            `json:"page"`
	Rows       int            `json:"rows"`
	TpResult   TypeResult     `json:"tp_result"`
	Unquoted   bool           `json:"unquoted"`
}

/**
* newField
* @param column *Column
* @return *Field
**/
func newField(column *Column) *Field {
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

/**
* describe
* @return et.Json
**/
func (s *Field) describe() et.Json {
	definition, err := json.Marshal(s)
	if err != nil {
		return et.Json{}
	}

	result := et.Json{}
	err = json.Unmarshal(definition, &result)
	if err != nil {
		return et.Json{}
	}

	return result
}

/**
* setValue
* @param value interface{}
**/
func (s *Field) setValue(value interface{}) {
	regexpMust := func(pattern string, value interface{}) (string, bool) {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(value.(string))

		if len(matches) > 1 {
			return matches[1], true
		} else {
			return value.(string), false
		}
	}

	switch value.(type) {
	case string:
		result, ok := regexpMust(`(?i)^CALC\((.*)\)$`, value)
		if ok {
			s.Value = result
			s.Unquoted = true
		} else {
			s.Value = value
		}
	default:
		s.Value = value
	}
}

/**
* setAgregation
* @param agr TypeAgregation
**/
func (s *Field) setAgregation(agr TypeAgregation) {
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
* ValueQuoted
* @return any
**/
func (s *Field) ValueQuoted() any {
	if s.Unquoted {
		return strs.Format(`%v`, s.Value)
	}

	return utility.Quote(s.Value)
}

/**
* ValueUnquoted
* @return any
**/
func (s *Field) ValueUnquoted() any {
	return utility.Unquote(s.Value)
}

/**
* asField
* @return string
**/
func (s *Field) asField() string {
	result := ""
	result = strs.Append(result, s.Schema, "")
	result = strs.Append(result, s.Table, ".")
	result = strs.Append(result, s.Source, ".")
	result = strs.Append(result, s.Name, ".")

	return result
}

/**
* asName
* @return string
**/
func (s *Field) asName() string {
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
	result := newField(s)

	return result
}
