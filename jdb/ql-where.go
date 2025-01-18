package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
)

type Connector int

const (
	NoC Connector = iota
	And
	Or
)

func (s Connector) Str() string {
	switch s {
	case And:
		return "and"
	case Or:
		return "or"
	default:
		return ""
	}
}

type Operator int

const (
	NoP Operator = iota
	Equal
	Neg
	In
	Like
	More
	Less
	MoreEq
	LessEq
	Between
	IsNull
	NotNull
	Search
)

func (s Operator) Str() string {
	switch s {
	case Equal:
		return "="
	case Neg:
		return "!="
	case In:
		return "in"
	case Like:
		return "like"
	case More:
		return ">"
	case Less:
		return "<"
	case MoreEq:
		return ">="
	case LessEq:
		return "<="
	case Between:
		return "between"
	case
		IsNull:
		return "is null"
	case NotNull:
		return "is not null"
	case Search:
		return "search"
	default:
		return ""
	}
}

type QlWhere struct {
	Conector Connector
	Key      interface{}
	Operator Operator
	Values   []interface{}
	Language string
}

/**
* NewQlWhere
* @params key interface{}
* @return QlWhere
**/
func NewQlWhere(key interface{}) *QlWhere {
	return &QlWhere{
		Conector: NoC,
		Key:      key,
		Operator: NoP,
		Values:   make([]interface{}, 0),
	}
}

/**
* getValue
* @param val interface{}
* @return string
**/
func (s *QlWhere) GetValue(val interface{}) string {
	switch v := val.(type) {
	case *QlSelect:
		return v.Field.AsField()
	case *Field:
		return v.AsField()
	case []interface{}:
		var result string
		for _, w := range v {
			val := s.GetValue(w)
			result = strs.Append(result, strs.Format(`%v`, val), ",")
		}
		return result
	default:
		return strs.Format(`%v`, utility.Quote(v))
	}
}

/**
* GetKey
* @return string
**/
func (s *QlWhere) GetKey() string {
	return s.GetValue(s.Key)
}

/**
* String
* @return string
**/
func (s *QlWhere) String() string {
	var result string

	if s.Conector != NoC {
		result = strs.Append(result, s.Conector.Str(), " ")
	}
	result = strs.Append(result, s.GetKey(), " ")
	result = strs.Append(result, s.Operator.Str(), " ")

	for _, val := range s.Values {
		result = strs.Append(result, s.GetValue(val), " ")
	}

	return result
}

/**
* Where
* @param val interface{}
* @return *QlFilter
**/
func (s *Ql) Where(val interface{}) *Ql {
	switch v := val.(type) {
	case string:
		field := s.GetField(v, false)
		if field != nil {
			s.where = NewQlWhere(field)
			return s
		}
	}

	s.where = NewQlWhere(val)
	return s
}

/**
* And
* @param val interface{}
* @return *QlFilter
**/
func (s *Ql) And(val interface{}) *QlFilter {
	result := s.Where(val)
	result.where.Conector = And

	return result.QlFilter
}

/**
* And
* @param val interface{}
* @return *QlFilter
**/
func (s *Ql) Or(val interface{}) *QlFilter {
	result := s.Where(val)
	result.where.Conector = Or

	return result.QlFilter
}

/**
* setWheres
* @param wheres []et.Json
**/
func (s *Ql) setWheres(wheres []et.Json) *Ql {
	for _, item := range wheres {
		if item["key"] != nil {
			s.Where(item["key"])
		} else if item["and"] != nil {
			s.And(item["and"])
		} else if item["or"] != nil {
			s.Or(item["or"])
		}

		s.setCondition(item)
	}

	return s
}

/**
* listWheres
* @return []string
**/
func (s *Ql) listWheres() []string {
	result := []string{}
	for _, val := range s.Wheres {
		result = append(result, val.String())
	}

	return result
}
