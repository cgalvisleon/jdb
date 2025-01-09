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

type LinqWhere struct {
	Conector Connector
	Key      interface{}
	Operator Operator
	Values   []interface{}
}

/**
* NewLinqWhere
* @params key interface{}
* @return LinqWhere
**/
func NewLinqWhere(key interface{}) *LinqWhere {
	return &LinqWhere{
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
func (s *LinqWhere) GetValue(val interface{}) string {
	switch v := val.(type) {
	case *LinqSelect:
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
func (s *LinqWhere) GetKey() string {
	return s.GetValue(s.Key)
}

/**
* String
* @return string
**/
func (s *LinqWhere) String() string {
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
* @return *LinqFilter
**/
func (s *Linq) Where(val interface{}) *Linq {
	field, ok := val.(string)
	if ok {
		field := s.GetField(field)
		if field != nil {
			s.where = NewLinqWhere(field)
			return s
		}
	}

	s.where = NewLinqWhere(val)
	return s
}

/**
* And
* @param val interface{}
* @return *LinqFilter
**/
func (s *Linq) And(val interface{}) *LinqFilter {
	result := s.Where(val)
	result.where.Conector = And

	return result.LinqFilter
}

/**
* And
* @param val interface{}
* @return *LinqFilter
**/
func (s *Linq) Or(val interface{}) *LinqFilter {
	result := s.Where(val)
	result.where.Conector = Or

	return result.LinqFilter
}

/**
* setWheres
* @param wheres []et.Json
**/
func (s *Linq) setWheres(wheres []et.Json) *Linq {
	for _, item := range wheres {
		if item["key"] != nil {
			s.Where(item["key"])
		} else if item["and"] != nil {
			s.And(item["and"])
		} else if item["or"] != nil {
			s.Or(item["or"])
		}

		if item["eq"] != nil {
			s.Eq(item["eq"])
		} else if item["neg"] != nil {
			s.Neg(item["neg"])
		} else if item["in"] != nil {
			s.In(item["in"])
		} else if item["like"] != nil {
			s.Like(item["like"])
		} else if item["more"] != nil {
			s.More(item["more"])
		} else if item["less"] != nil {
			s.Less(item["less"])
		} else if item["moreEq"] != nil {
			s.MoreEq(item["moreEq"])
		} else if item["lessEq"] != nil {
			s.LessEs(item["lessEq"])
		} else if item["search"] != nil {
			s.Search(item["search"])
		} else if item["between"] != nil {
			s.Between(item["between"])
		} else if item["isNull"] != nil {
			s.IsNull()
		} else if item["notNull"] != nil {
			s.NotNull()
		}
	}

	return s
}

/**
* listWheres
* @return []string
**/
func (s *Linq) listWheres() []string {
	result := []string{}
	for _, val := range s.Wheres {
		result = append(result, val.String())
	}

	return result
}
