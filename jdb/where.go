package jdb

import (
	"fmt"
	"slices"
	"strings"

	"github.com/cgalvisleon/et/et"
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
	case IsNull:
		return "is null"
	case NotNull:
		return "is not null"
	case Search:
		return "search"
	default:
		return "Any"
	}
}

/**
* Name
* @return string
**/
func (s *Operator) Name() string {
	switch *s {
	case Equal:
		return "Equal"
	case Neg:
		return "Neg"
	case In:
		return "In"
	case Like:
		return "Like"
	case More:
		return "More"
	case Less:
		return "Less"
	case MoreEq:
		return "MoreEq"
	case LessEq:
		return "LessEq"
	case Between:
		return "Between"
	case IsNull:
		return "IsNull"
	case NotNull:
		return "NotNull"
	case Search:
		return "Search"
	default:
		return "Any"
	}
}

/**
* Command
* @return string
**/
func (s *Operator) command() string {
	switch *s {
	case Equal:
		return "eq"
	case Neg:
		return "neg"
	case In:
		return "in"
	case Like:
		return "like"
	case More:
		return "more"
	case Less:
		return "less"
	case MoreEq:
		return "moreEq"
	case LessEq:
		return "lessEq"
	case Between:
		return "between"
	case IsNull:
		return "isNull"
	case NotNull:
		return "notNull"
	case Search:
		return "search"
	default:
		return "any"
	}
}

/**
* OperatorToCommand
* @param op Operator
* @return string
**/
func OperatorToCommand(op Operator) string {
	return op.command()
}

/**
* StrToOperator
* @param str string
* @return Operator
**/
func StrToOperator(str string) Operator {
	switch str {
	case "eq":
		return Equal
	case "neg":
		return Neg
	case "in":
		return In
	case "like":
		return Like
	case "more":
		return More
	case "less":
		return Less
	case "moreEq":
		return MoreEq
	case "lessEq":
		return LessEq
	case "between":
		return Between
	case "isNull":
		return IsNull
	case "notNull":
		return NotNull
	case "search":
		return Search
	default:
		return NoP
	}
}

type QlCondition struct {
	Connector Connector         `json:"connector"`
	Field     interface{}       `json:"field"`
	Operator  Operator          `json:"operator"`
	Value     interface{}       `json:"value"`
	Language  string            `json:"language"`
	validator func(val any) any `json:"-"`
}

/**
* newQlCondition
* @params field interface{}, validator func(val any) any
* @return QlWhere
**/
func newQlCondition(field interface{}, validator func(val any) any) *QlCondition {
	return &QlCondition{
		Connector: NoC,
		Field:     field,
		Operator:  NoP,
		Value:     []interface{}{},
		validator: validator,
	}
}

/**
* setVal
* @param val interface{}
* @return *QlCondition
**/
func (s *QlCondition) setVal(val interface{}) {
	val = s.validator(val)
	switch v := val.(type) {
	case *Field:
		s.Value = v
	case Field:
		s.Value = v
	case Column:
		s.Value = GetField(&v)
	case *Column:
		s.Value = GetField(v)
	default:
		s.Value = val
	}
}

/**
* getField
* @return string
**/
func (s *QlCondition) getField() string {
	switch v := s.Field.(type) {
	case *Field:
		return v.asName()
	case Field:
		return v.asName()
	default:
		return fmt.Sprintf(`%v`, Quote(v))
	}
}

/**
* getValue
* @return string
**/
func (s *QlCondition) getValue() string {
	switch v := s.Value.(type) {
	case *Field:
		return v.asName()
	case Field:
		return v.asName()
	default:
		return fmt.Sprintf(`%v`, v)
	}
}

/**
* getCondition
* @return et.Json
**/
func (s *QlCondition) getCondition() et.Json {
	return et.Json{
		s.Operator.command(): s.getValue(),
	}
}

type QlWhere struct {
	Wheres    []*QlCondition    `json:"wheres"`
	IsDebug   bool              `json:"-"`
	language  string            `json:"-"`
	validator func(val any) any `json:"-"`
}

/**
* newQlWhere
* @return *QlWhere
**/
func newQlWhere(validator func(val any) any) *QlWhere {
	return &QlWhere{
		Wheres:    []*QlCondition{},
		IsDebug:   false,
		validator: validator,
	}
}

/**
* setDebug
* @param debug bool
* @return *Ql
**/
func (s *QlWhere) setDebug(debug bool) *QlWhere {
	s.IsDebug = debug

	return s
}

/**
* getWheres
* @return et.Json
**/
func (s *QlWhere) getWheres() et.Json {
	result := et.Json{}
	and := []et.Json{}
	or := []et.Json{}
	for i, condition := range s.Wheres {
		if condition.Field == nil {
			continue
		}

		field := condition.getField()
		if condition.Connector == And {
			and = append(and, et.Json{field: condition.getCondition()})
		} else if condition.Connector == Or {
			or = append(or, et.Json{field: condition.getCondition()})
		} else if i == 0 {
			result.Set(field, condition.getCondition())
		}
	}

	if len(and) > 0 {
		result.Set("AND", and)
	}
	if len(or) > 0 {
		result.Set("OR", or)
	}

	return result
}

/**
* setWhere
* @param val field interface{}
* @return *QlWhere
**/
func (s *QlWhere) setWhere(field interface{}) *QlWhere {
	where := newQlCondition(field, s.validator)
	if len(s.Wheres) > 0 {
		where.Connector = And
	}

	s.Wheres = append(s.Wheres, where)

	return s
}

/**
* setAnd
* @param val field interface{}
* @return *QlWhere
**/
func (s *QlWhere) setAnd(field interface{}) *QlWhere {
	return s.setWhere(field)
}

/**
* setOr
* @param val field interface{}
* @return *QlWhere
**/
func (s *QlWhere) setOr(field interface{}) *QlWhere {
	where := newQlCondition(field, s.validator)
	if len(s.Wheres) > 0 {
		where.Connector = Or
	}

	s.Wheres = append(s.Wheres, where)

	return s
}

/**
* setValue
* @param values et.Json
* @return *QlWhere
**/
func (s *QlWhere) setValue(values et.Json) *QlWhere {
	for key, val := range values {
		val = s.validator(val)
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
* Where
* @param val interface{}
* @return *QlWhere
**/
func (s *QlWhere) Where(val interface{}) *QlWhere {
	val = s.validator(val)
	if val != nil {
		s.setWhere(val)
	}

	return s
}

/**
* And
* @param val interface{}
* @return *QlWhere
**/
func (s *QlWhere) And(val interface{}) *QlWhere {
	return s.Where(val)
}

/**
* Or
* @param val interface{}
* @return *QlWhere
**/
func (s *QlWhere) Or(val interface{}) *QlWhere {
	val = s.validator(val)
	if val != nil {
		if len(s.Wheres) == 0 {
			s.setWhere(val)
		} else {
			s.setOr(val)
		}
	}

	return s
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
* SetWheres
* @param wheres et.Json
* @return *QlWhere
**/
func (s *QlWhere) SetWheres(wheres et.Json) *QlWhere {
	and := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				s.And(key).setValue(val.Json(key))
			}
		}
	}

	or := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				s.Or(key).setValue(val.Json(key))
			}
		}
	}

	for key := range wheres {
		key = strings.ToLower(key)
		if slices.Contains([]string{"and", "or"}, key) {
			continue
		}

		val := wheres.Json(key)
		s.Where(key).setValue(val)
	}

	for key := range wheres {
		switch strings.ToLower(key) {
		case "and":
			vals := wheres.ArrayJson(key)
			and(vals)
		case "or":
			vals := wheres.ArrayJson(key)
			or(vals)
		}
	}

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
