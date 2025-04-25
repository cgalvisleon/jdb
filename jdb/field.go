package jdb

import (
	"encoding/json"
	"regexp"
	"strconv"
	"time"

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

/**
* Describe
* @return et.Json
**/
func (s *Field) Describe() et.Json {
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
* SetValue
* @param value interface{}
**/
func (s *Field) SetValue(value interface{}) {
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
* SetAgregation
* @param agr TypeAgregation
**/
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
* AsString
* @return string
**/
func (s *Field) AsStr() string {
	if s.Value == nil {
		return ""
	}

	return strs.Format("%v", s.Value)
}

/**
* AsInt
* @return int
**/
func (s *Field) AsInt() int {
	if s.Value == nil {
		return 0
	}

	result, err := strconv.Atoi(s.Value.(string))
	if err != nil {
		return 0
	}

	return result
}

/**
* AsInt64
* @return int64
**/
func (s *Field) AsInt64() int64 {
	if s.Value == nil {
		return 0
	}

	result, err := strconv.ParseInt(s.Value.(string), 10, 64)
	if err != nil {
		return 0
	}

	return result
}

/**
* AsFloat
* @return float64
**/
func (s *Field) AsFloat() float64 {
	if s.Value == nil {
		return 0
	}

	result, err := strconv.ParseFloat(s.Value.(string), 64)
	if err != nil {
		return 0
	}

	return result
}

/**
* AsBool
* @return bool
**/
func (s *Field) AsBool() bool {
	if s.Value == nil {
		return false
	}

	result, err := strconv.ParseBool(s.Value.(string))
	if err != nil {
		return false
	}

	return result
}

/**
* AsTime
* @return time.Time
**/
func (s *Field) AsTime() time.Time {
	if s.Value == nil {
		return time.Time{}
	}

	result, err := time.Parse(time.RFC3339, s.Value.(string))
	if err != nil {
		return time.Time{}
	}

	return result
}

/**
* AsJson
* @return et.Json
**/
func (s *Field) AsJson() et.Json {
	if s.Value == nil {
		return et.Json{}
	}

	bt, err := json.Marshal(s.Value)
	if err != nil {
		return et.Json{}
	}

	var result et.Json
	err = json.Unmarshal(bt, &result)
	if err != nil {
		return et.Json{}
	}

	return result
}

/**
* GetField
* @return *Field
**/
func (s *Column) GetField() *Field {
	result := NewField(s)

	return result
}
