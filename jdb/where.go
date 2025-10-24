package jdb

import (
	"fmt"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type Field struct {
	Value string `json:"value"`
}

/**
* F
* @param val string
* @return *Field
**/
func F(val string) *Field {
	return &Field{
		Value: val,
	}
}

type Condition struct {
	Field string      `json:"field"`
	Op    string      `json:"op"`
	Value interface{} `json:"value"`
}

/**
* ToJson
* @return et.Json
**/
func (s *Condition) ToJson() et.Json {
	return et.Json{
		s.Field: et.Json{
			s.Op: s.Value,
		},
	}
}

/**
* condition
* @param field string, value interface{}, op string
* @return *Condition
**/
func condition(field string, value interface{}, op string) *Condition {
	switch v := value.(type) {
	case *Field:
		value = v.Value
	default:
		value = Quote(value)
	}

	return &Condition{
		Field: field,
		Op:    op,
		Value: value,
	}
}

/**
* Eq
* @param field string, value interface{}
* @return Condition
**/
func Eq(field string, value interface{}) *Condition {
	return condition(field, value, "eq")
}

/**
* Neg
* @param field string, value interface{}
* @return Condition
**/
func Neg(field string, value interface{}) *Condition {
	return condition(field, value, "ne")
}

/**
* Less
* @param field string, value interface{}
* @return Condition
**/
func Less(field string, value interface{}) *Condition {
	return condition(field, value, "less")
}

/**
* LessEq
* @param field string, value interface{}
* @return Condition
**/
func LessEq(field string, value interface{}) *Condition {
	return condition(field, value, "less_eq")
}

/**
* More
* @param field string, value interface{}
* @return Condition
**/
func More(field string, value interface{}) *Condition {
	return condition(field, value, "more")
}

/**
* MoreEq
* @param field string, value interface{}
* @return Condition
**/
func MoreEq(field string, value interface{}) *Condition {
	return condition(field, value, "more_eq")
}

/**
* Like
* @param field string, value interface{}
* @return Condition
**/
func Like(field string, value interface{}) *Condition {
	return condition(field, value, "like")
}

/**
* Ilike
* @param field string, value interface{}
* @return Condition
**/
func Ilike(field string, value interface{}) *Condition {
	return condition(field, value, "ilike")
}

/**
* In
* @param field string, value interface{}
* @return Condition
**/
func In(field string, value interface{}) *Condition {
	return condition(field, value, "in")
}

/**
* NotIn
* @param field string, value interface{}
* @return Condition
**/
func NotIn(field string, value interface{}) *Condition {
	return condition(field, value, "not_in")
}

/**
* Is
* @param field string, value interface{}
* @return Condition
**/
func Is(field string, value interface{}) *Condition {
	return condition(field, value, "is")
}

/**
* IsNot
* @param field string, value interface{}
* @return Condition
**/
func IsNot(field string, value interface{}) *Condition {
	return condition(field, value, "is_not")
}

/**
* Null
* @param field string
* @return Condition
**/
func Null(field string) *Condition {
	return condition(field, nil, "null")
}

/**
* NotNull
* @param field string
* @return Condition
**/
func NotNull(field string) *Condition {
	return condition(field, nil, "not_null")
}

/**
* Between
* @param field string, value []interface{}
* @return Condition
**/
func Between(field string, value []interface{}) *Condition {
	return condition(field, value, "between")
}

/**
* NotBetween
* @param field string, value []interface{}
* @return Condition
**/
func NotBetween(field string, value []interface{}) *Condition {
	return condition(field, value, "not_between")
}

type where struct {
	Wheres []et.Json `json:"where"`
	From   *Model    `json:"from"`
	As     string    `json:"as"`
}

/**
* newWhere
* @param model *Model, as string
* @return *where
**/
func newWhere(model *Model, as string) *where {
	return &where{
		Wheres: []et.Json{},
		From:   model,
		As:     as,
	}
}

/**
* where
* @param cond Condition
* @return *where
**/
func (s *where) where(cond *Condition, conector string) *where {
	if s.From != nil {
		f := cond.Field
		col, ok := s.From.GetColumn(f)
		if !ok {
			if s.From.UseAtribs() {
				f = GetFieldName(f)
				cond.Field = fmt.Sprintf("%s:%s", s.From.SourceField, f)
				cond.Field = strs.Append(s.As, cond.Field, ".")
			}
		} else {
			tp := col.String("type")
			if TypeColumn[tp] {
				cond.Field = f
				cond.Field = strs.Append(s.As, cond.Field, ".")
			}
			if TypeAtrib[tp] {
				f = GetFieldName(f)
				cond.Field = fmt.Sprintf("%s:%s", s.From.SourceField, f)
				cond.Field = strs.Append(s.As, cond.Field, ".")
			}
		}
	}

	if len(s.Wheres) == 0 {
		s.Wheres = append(s.Wheres, et.Json{
			cond.Field: et.Json{
				cond.Op: cond.Value,
			},
		})
	} else {
		conds := []et.Json{}
		idx := slices.IndexFunc(s.Wheres, func(v et.Json) bool { return strs.Uppcase(v.String(conector)) == strs.Uppcase(conector) })
		if idx != -1 {
			conds = s.Wheres[idx].ArrayJson(conector)
		}

		conds = append(conds, et.Json{
			cond.Field: et.Json{
				cond.Op: cond.Value,
			},
		})

		if idx == -1 {
			s.Wheres = append(s.Wheres, et.Json{
				conector: conds,
			})
		} else {
			s.Wheres[idx][conector] = conds
		}
	}

	return s
}

/**
* Where
* @param cond Condition
* @return *where
**/
func (s *where) Where(cond *Condition) *where {
	return s.where(cond, "and")
}

/**
* And
* @param cond Condition
* @return *where
**/
func (s *where) And(cond *Condition) *where {
	return s.where(cond, "and")
}

/**
* Or
* @param cond Condition
* @return *where
**/
func (s *where) Or(cond *Condition) *where {
	return s.where(cond, "or")
}
