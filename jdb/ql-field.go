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
	ExtractYear
	ExtractMonth
	ExtractDay
	ExtractHour
	ExtractMinute
	ExtractSecond
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
	case ExtractYear:
		return "YEAR"
	case ExtractMonth:
		return "MONTH"
	case ExtractDay:
		return "DAY"
	case ExtractHour:
		return "HOUR"
	case ExtractMinute:
		return "MINUTE"
	case ExtractSecond:
		return "SECOND"
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
	Column     *Column        `json:"-"`
	Schema     string         `json:"schema"`
	Model      string         `json:"model"`
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
	Select     []interface{}  `json:"select"`
	Joins      []et.Json      `json:"joins"`
	Where      et.Json        `json:"where"`
	GroupBy    []string       `json:"group_by"`
	Havings    et.Json        `json:"havings"`
	OrderBy    et.Json        `json:"order_by"`
}

func (s *Field) describe() et.Json {
	return et.Json{
		"column_type": s.Column.TypeColumn.Str(),
		"schema":      s.Schema,
		"model":       s.Model,
		"as":          s.As,
		"name":        s.Name,
		"source":      s.Source,
		"agregation":  s.Agregation.Str(),
		"value":       s.Value,
		"alias":       s.Alias,
		"hidden":      s.Hidden,
		"page":        s.Page,
		"rows":        s.Rows,
		"tp_result":   s.TpResult,
		"unquoted":    s.Unquoted,
		"select":      s.Select,
		"joins":       s.Joins,
		"where":       s.Where,
		"group_by":    s.GroupBy,
		"havings":     s.Havings,
		"order_by":    s.OrderBy,
	}
}

/**
* newField
* @param name string
* @return *Field
**/
func newField(name string) *Field {
	return &Field{
		Name:    name,
		Select:  make([]interface{}, 0),
		Joins:   make([]et.Json, 0),
		Where:   et.Json{},
		GroupBy: make([]string, 0),
		Havings: et.Json{},
		OrderBy: et.Json{},
	}
}

/**
* Serialize
* @return []byte, error
**/
func (s *Field) Serialize() ([]byte, error) {
	result, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}

	return result, nil
}

/**
* Describe
* @return et.Json
**/
func (s *Field) Describe() et.Json {
	definition, err := s.Serialize()
	if err != nil {
		return et.Json{}
	}

	result := et.Json{}
	err = json.Unmarshal(definition, &result)
	if err != nil {
		return et.Json{}
	}

	result["column"] = s.Column.Describe()

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
* asField
* @return string
**/
func (s *Field) asField() string {
	result := ""
	result = strs.Append(result, s.Schema, "")
	result = strs.Append(result, s.Model, ".")
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
		return strs.Format(`%s.%s`, s.As, s.Name)
	}

	return strs.Format(`%v`, s.Name)
}

/**
* GetField
* @return *Field
**/
func (s *Column) GetField() *Field {
	result := newField(s.Name)
	result.Column = s
	result.Name = s.Name
	result.Alias = s.Name
	result.Hidden = s.Hidden

	if s.TypeColumn == TpRelatedTo {
		result.Page = 1
		result.Rows = 30
		result.TpResult = TpResult
	}

	if s.Model == nil {
		return result
	}

	result.Schema = s.Model.Schema
	result.Model = s.Model.Name

	if s.Source == nil {
		return result
	}

	result.Source = s.Source.Name

	return result
}
