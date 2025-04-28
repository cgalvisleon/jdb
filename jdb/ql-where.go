package jdb

import (
	"fmt"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type QlWhere struct {
	Wheres   []*QlCondition
	getField func(name string) *Field
	history  bool
	language string
	IsDebug  bool
}

/**
* NewQlWhere
* @return *QlWhere
**/
func NewQlWhere() *QlWhere {
	return &QlWhere{
		Wheres:  []*QlCondition{},
		history: false,
		IsDebug: false,
	}
}

/**
* where
* @param val field interface{}
* @return *QlWhere
**/
func (s *QlWhere) where(field interface{}) *QlWhere {
	where := NewQlCondition(field)
	s.Wheres = append(s.Wheres, where)

	return s
}

/**
* and
* @param val field interface{}
* @return *QlWhere
**/
func (s *QlWhere) and(field interface{}) *QlWhere {
	where := NewQlCondition(field)
	where.Connector = And
	s.Wheres = append(s.Wheres, where)

	return s
}

/**
* or
* @param val field interface{}
* @return *QlWhere
**/
func (s *QlWhere) or(field interface{}) *QlWhere {
	where := NewQlCondition(field)
	where.Connector = Or
	s.Wheres = append(s.Wheres, where)

	return s
}

/**
* condition
* @return *QlCondition
**/
func (s *QlWhere) condition() *QlCondition {
	idx := len(s.Wheres)
	if idx <= 0 {
		return nil
	}

	return s.Wheres[idx-1]
}

/**
* String
* @return string
**/
func (s *QlWhere) string() string {
	var result string
	for _, val := range s.Wheres {
		result = strs.Append(result, val.string(), " ")
	}

	return result
}

/**
* Eq
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Eq(val interface{}) *QlWhere {
	condition := s.condition()
	if condition == nil {
		return s
	}

	condition.Operator = Equal
	condition.setVal(val)

	return s
}

/**
* Neg
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Neg(val interface{}) *QlWhere {
	condition := s.condition()
	if condition == nil {
		return s
	}

	condition.Operator = Neg
	condition.setVal(val)

	return s
}

/**
* In
* @param val ...any
* @return QlWhere
**/
func (s *QlWhere) In(val ...any) *QlWhere {
	condition := s.condition()
	if condition == nil {
		return s
	}

	condition.Operator = In
	condition.setVal(val)

	return s
}

/**
* Like
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Like(val interface{}) *QlWhere {
	condition := s.condition()
	if condition == nil {
		return s
	}

	condition.Operator = Like
	condition.setVal(val)

	return s
}

/**
* More
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) More(val interface{}) *QlWhere {
	condition := s.condition()
	if condition == nil {
		return s
	}

	condition.Operator = More
	condition.setVal(val)

	return s
}

/**
* Less
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Less(val interface{}) *QlWhere {
	condition := s.condition()
	if condition == nil {
		return s
	}

	condition.Operator = Less
	condition.setVal(val)

	return s
}

/**
* MoreEq
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) MoreEq(val interface{}) *QlWhere {
	condition := s.condition()
	if condition == nil {
		return s
	}

	condition.Operator = MoreEq
	condition.setVal(val)

	return s
}

/**
* LessEq
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) LessEq(val interface{}) *QlWhere {
	condition := s.condition()
	if condition == nil {
		return s
	}

	condition.Operator = LessEq
	condition.setVal(val)

	return s
}

/**
* Between
* @param val1, val2 interface{}
* @return QlWhere
**/
func (s *QlWhere) Between(vals interface{}) *QlWhere {
	val, ok := vals.([]interface{})
	if !ok {
		return s
	}

	condition := s.condition()
	if condition == nil {
		return s
	}

	condition.Operator = Between
	condition.setVal(val)

	return s
}

/**
* Search
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Search(language string, val interface{}) *QlWhere {
	condition := s.condition()
	if condition == nil {
		return s
	}

	condition.Operator = Search
	condition.Language = language
	condition.setVal(val)

	return s
}

/**
* IsNull
* @return *QlWhere
**/
func (s *QlWhere) IsNull() *QlWhere {
	condition := s.condition()
	if condition == nil {
		return s
	}

	condition.Operator = IsNull

	return s
}

/**
* NotNull
* @return *QlWhere
**/
func (s *QlWhere) NotNull() *QlWhere {
	condition := s.condition()
	if condition == nil {
		return s
	}

	condition.Operator = NotNull

	return s
}

/**
* History
* @param v bool
* @return *Command
**/
func (s *QlWhere) History(v bool) *QlWhere {
	s.history = v

	return s
}

/**
* Debug
* @param v bool
* @return *Command
**/
func (s *QlWhere) Debug() *QlWhere {
	s.IsDebug = true

	return s
}

/**
* setValue
* @param values et.Json
* @return *QlWhere
**/
func (s *QlWhere) setValue(values et.Json, validator func(val interface{}) interface{}) *QlWhere {
	for key, val := range values {
		val = validator(val)
		switch key {
		case "eq":
			s.Eq(val)
		case "neg":
			s.Neg(val)
		case "in":
			s.In(val)
		case "like":
			s.Like(val)
		case "more":
			s.More(val)
		case "less":
			s.Less(val)
		case "moreEq":
			s.MoreEq(val)
		case "lessEq":
			s.LessEq(val)
		case "between":
			s.Between(val)
		case "isNull":
			s.IsNull()
		case "notNull":
			s.NotNull()
		case "search":
			s.Search(s.language, val)
		}
	}

	return s
}

/**
* listWheres
* @param asField func(field *Field) string
* @return et.Json
**/
func (s *QlWhere) listWheres() et.Json {
	result := et.Json{}
	and := []et.Json{}
	or := []et.Json{}
	for i, condition := range s.Wheres {
		if condition.Field == nil {
			continue
		}

		field := fmt.Sprintf(`%v`, condition.getField())
		def := et.Json{condition.Operator.Str(): condition.getValue()}
		if condition.Connector == And {
			and = append(and, et.Json{field: def})
		} else if condition.Connector == Or {
			or = append(or, et.Json{field: def})
		} else if i == 0 {
			result.Set(field, def)
		}
	}

	if len(and) > 0 {
		result.Set("and", and)
	}
	if len(or) > 0 {
		result.Set("or", or)
	}

	return result
}

/**
* Where
* @param val interface{}
* @return *Ql
**/
func (s *Ql) Where(val interface{}) *Ql {
	val = s.validator(val)
	if val != nil {
		if len(s.Wheres) == 0 {
			s.where(val)
		} else {
			s.and(val)
		}
	}

	return s
}

/**
* And
* @param val interface{}
* @return *Ql
**/
func (s *Ql) And(val interface{}) *Ql {
	val = s.validator(val)
	if val != nil {
		s.and(val)
	}

	return s
}

/**
* Or
* @param val interface{}
* @return *Ql
**/
func (s *Ql) Or(val interface{}) *Ql {
	val = s.validator(val)
	if val != nil {
		s.or(val)
	}

	return s
}

/**
* Eq
* @param val interface{}
* @return *Ql
**/
func (s *Ql) Eq(val interface{}) *Ql {
	s.QlWhere.Eq(val)

	return s
}

/**
* Neg
* @param val interface{}
* @return *Ql
**/
func (s *Ql) Neg(val interface{}) *Ql {
	s.QlWhere.Neg(val)

	return s
}

/**
* In
* @param val ...any
* @return *Ql
**/
func (s *Ql) In(val ...any) *Ql {
	s.QlWhere.In(val...)

	return s
}

/**
* Like
* @param val interface{}
* @return *Ql
**/
func (s *Ql) Like(val interface{}) *Ql {
	s.QlWhere.Like(val)

	return s
}

/**
* More
* @param val interface{}
* @return *Ql
**/
func (s *Ql) More(val interface{}) *Ql {
	s.QlWhere.More(val)

	return s
}

/**
* Less
* @param val interface{}
* @return *Ql
**/
func (s *Ql) Less(val interface{}) *Ql {
	s.QlWhere.Less(val)

	return s
}

/**
* MoreEq
* @param val interface{}
* @return *Ql
**/
func (s *Ql) MoreEq(val interface{}) *Ql {
	s.QlWhere.MoreEq(val)

	return s
}

/**
* LessEq
* @param val interface{}
* @return *Ql
**/
func (s *Ql) LessEq(val interface{}) *Ql {
	s.QlWhere.LessEq(val)

	return s
}

/*
*
* Between
* @param vals interface{}
* @return *Ql
**/
func (s *Ql) Between(vals interface{}) *Ql {
	s.QlWhere.Between(vals)

	return s
}

/**
* Search
* @param language string, val interface{}
* @return *Ql
**/
func (s *Ql) Search(language string, val interface{}) *Ql {
	s.QlWhere.Search(language, val)

	return s
}

/**
* IsNull
* @return *Ql
**/
func (s *Ql) IsNull() *Ql {
	s.QlWhere.IsNull()

	return s
}

/**
* NotNull
* @return *Ql
**/
func (s *Ql) NotNull() *Ql {
	s.QlWhere.NotNull()

	return s
}

/**
* History
* @param v bool
* @return *Ql
**/
func (s *Ql) History(v bool) *Ql {
	s.QlWhere.History(v)

	return s
}

/**
* Debug
* @param v bool
* @return *Ql
**/
func (s *Ql) Debug() *Ql {
	s.QlWhere.Debug()

	return s
}

/**
* setDebug
* @param debug bool
* @return *Ql
**/
func (s *Ql) setDebug(debug bool) *Ql {
	s.IsDebug = debug

	return s
}

/**
* setWheres
* @param wheres et.Json
* @return *Ql
**/
func (s *Ql) setWheres(wheres et.Json) *Ql {
	and := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				s.And(key).setValue(val.Json(key), s.validator)
			}
		}
	}

	or := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				s.Or(key).setValue(val.Json(key), s.validator)
			}
		}
	}

	for key := range wheres {
		if slices.Contains([]string{"and", "AND", "or", "OR"}, key) {
			continue
		}

		s.Where(key).setValue(wheres.Json(key), s.validator)
	}

	for key := range wheres {
		switch key {
		case "and", "AND":
			vals := wheres.ArrayJson(key)
			and(vals)
		case "or", "OR":
			vals := wheres.ArrayJson(key)
			or(vals)
		}
	}

	return s
}

/**
* listWheres
* @return et.Json
**/
func (s *Ql) listWheres() et.Json {
	return s.QlWhere.listWheres()
}
