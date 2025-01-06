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
	Value    []interface{}
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

	for _, val := range s.Value {
		result = strs.Append(result, s.GetValue(val), " ")
	}

	return result
}

type FilterTo interface {
	And(val interface{}) *LinqFilter
	Or(fval interface{}) *LinqFilter
	Select(fields ...string) *Linq
	Data(fields ...string) *Linq
	Return(fields ...string) *Command
}

type LinqFilter struct {
	main  FilterTo
	where *LinqWhere
}

/**
* NewLinqFilter
* @param main interface{}
* @param key interface{}
* @return *LinqFilter
**/
func NewLinqFilter(main FilterTo, key interface{}) *LinqFilter {
	where := &LinqWhere{
		Conector: NoC,
		Key:      key,
		Operator: NoP,
		Value:    make([]interface{}, 0),
	}

	return &LinqFilter{
		main:  main,
		where: where,
	}
}

/**
* Add
* @param val interface{}
* @return *LinqFilter
**/
func (s *LinqFilter) Add(val interface{}) FilterTo {
	appendValue := func(linq *Linq, value interface{}) {
		switch v := value.(type) {
		case string:
			field := linq.GetField(v)
			if field != nil {
				s.where.Value = append(s.where.Value, field)
			} else {
				s.where.Value = append(s.where.Value, value)
			}
		default:
			s.where.Value = append(s.where.Value, value)
		}
	}

	switch m := s.main.(type) {
	case *Linq:
		appendValue(m, val)
		m.Wheres = append(m.Wheres, s.where)
	case *Command:
		s.where.Value = append(s.where.Value, val)
		m.Wheres = append(m.Wheres, s.where)
	case *LinqJoin:
		appendValue(m.Linq, val)
		m.Wheres = append(m.Wheres, s.where)
	}

	return s.main
}

/**
* Eq
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) Eq(val interface{}) FilterTo {
	s.where.Operator = Equal
	return s.Add(val)
}

/**
* Neg
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) Neg(val interface{}) FilterTo {
	s.where.Operator = Neg
	return s.Add(val)
}

/**
* In
* @param val ...any
* @return FilterTo
**/
func (s *LinqFilter) In(val ...any) FilterTo {
	s.where.Operator = In
	return s.Add(val)
}

/**
* Like
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) Like(val interface{}) FilterTo {
	s.where.Operator = Like
	return s.Add(val)
}

/**
* More
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) More(val interface{}) FilterTo {
	s.where.Operator = More
	return s.Add(val)
}

/**
* Less
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) Less(val interface{}) FilterTo {
	s.where.Operator = Less
	return s.Add(val)
}

/**
* MoreEq
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) MoreEq(val interface{}) FilterTo {
	s.where.Operator = MoreEq
	return s.Add(val)
}

/**
* LessEq
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) LessEs(val interface{}) FilterTo {
	s.where.Operator = LessEq
	return s.Add(val)
}

/**
* Search
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) Search(val interface{}) FilterTo {
	s.where.Operator = Search
	return s.Add(val)
}

/**
* Between
* @param val1, val2 interface{}
* @return FilterTo
**/
func (s *LinqFilter) Between(val interface{}) FilterTo {
	s.where.Operator = Between
	vals, ok := val.([]interface{})
	if !ok {
		return s.main
	}

	switch len(vals) {
	case 1:
		return s.Add(vals[0])
	case 2:
		s.Add(vals[0])
		return s.Add(vals[1])
	default:
		return s.main
	}
}

/**
* IsNull
* @return *LinqFilter
**/
func (s *LinqFilter) IsNull() *LinqFilter {
	s.where.Operator = IsNull
	return s
}

/**
* NotNull
* @return *LinqFilter
**/
func (s *LinqFilter) NotNull() *LinqFilter {
	s.where.Operator = NotNull
	return s
}

/**
* Where
* @param val interface{}
* @return *LinqFilter
**/
func (s *Linq) Where(val interface{}) *LinqFilter {
	switch v := val.(type) {
	case string:
		sel := s.GetField(v)
		if sel != nil {
			return NewLinqFilter(s, sel)
		}
	}

	return NewLinqFilter(s, val)
}

/**
* And
* @param val interface{}
* @return *LinqFilter
**/
func (s *Linq) And(val interface{}) *LinqFilter {
	result := s.Where(val)
	result.where.Conector = And
	return result
}

/**
* And
* @param val interface{}
* @return *LinqFilter
**/
func (s *Linq) Or(val interface{}) *LinqFilter {
	result := s.Where(val)
	result.where.Conector = Or
	return result
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

/**
* setWheres
* @param query []et.Json
**/
func (s *Linq) setWheres(query []et.Json) *Linq {
	var filter *LinqFilter
	var filterTo FilterTo
	for _, item := range query {
		if item["key"] != nil {
			filter = s.Where(item["key"])
		} else if filterTo == nil {
			continue
		} else if item["and"] != nil {
			filter = filterTo.And(item["and"])
		} else if item["or"] != nil {
			filter = filterTo.Or(item["or"])
		}

		if filter == nil {
			continue
		} else if item["eq"] != nil {
			filterTo = filter.Eq(item["eq"])
		} else if item["neg"] != nil {
			filterTo = filter.Neg(item["neg"])
		} else if item["in"] != nil {
			filterTo = filter.In(item["in"])
		} else if item["like"] != nil {
			filterTo = filter.Like(item["like"])
		} else if item["more"] != nil {
			filterTo = filter.More(item["more"])
		} else if item["less"] != nil {
			filterTo = filter.Less(item["less"])
		} else if item["moreEq"] != nil {
			filterTo = filter.MoreEq(item["moreEq"])
		} else if item["lessEq"] != nil {
			filterTo = filter.LessEs(item["lessEq"])
		} else if item["search"] != nil {
			filterTo = filter.Search(item["search"])
		} else if item["between"] != nil {
			filterTo = filter.Between(item["between"])
		} else if item["isNull"] != nil {
			filter.IsNull()
		} else if item["notNull"] != nil {
			filter.NotNull()
		}

	}

	return s
}
