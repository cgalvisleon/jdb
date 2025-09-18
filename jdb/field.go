package jdb

import (
	"encoding/json"
	"fmt"
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

	switch v := value.(type) {
	case string:
		result, ok := regexpMust(`(?i)^CALC\((.*)\)$`, v)
		if ok {
			s.Value = result
			s.Unquoted = true
		} else {
			re := regexp.MustCompile(`^:(.*)`)
			matches := re.FindStringSubmatch(v)
			if len(matches) > 1 {
				s.Value = matches[1]
				s.Unquoted = true
			} else {
				s.Value = v
			}
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
		s.Alias = fmt.Sprintf("sum_%s", s.Name)
	case AgregationCount:
		s.Alias = fmt.Sprintf("count_%s", s.Name)
	case AgregationAvg:
		s.Alias = fmt.Sprintf("avg_%s", s.Name)
	case AgregationMin:
		s.Alias = fmt.Sprintf("min_%s", s.Name)
	case AgregationMax:
		s.Alias = fmt.Sprintf("max_%s", s.Name)
	}
}

/**
* ValueQuoted
* @return any
**/
func (s *Field) ValueQuoted() any {
	if s.Unquoted {
		return fmt.Sprintf(`%v`, s.Value)
	}

	return Quote(s.Value)
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
		return fmt.Sprintf(`%s.%s`, s.As, s.Name)
	}

	return fmt.Sprintf(`%v`, s.Name)
}

/**
* GetField
* @return *Field
**/
func GetField(col *Column) *Field {
	result := newField(col.Name)
	result.Column = col
	result.Name = col.Name
	result.Alias = col.Name
	result.Hidden = col.Hidden

	if col.TypeColumn == TpRelatedTo {
		result.Page = 1
		result.Rows = 30
		result.TpResult = TpResult
	}

	if col.Model == nil {
		return result
	}

	result.Schema = col.Model.Schema
	result.Model = col.Model.Name

	if col.Source == nil {
		return result
	}

	result.Source = col.Source.Name

	return result
}
