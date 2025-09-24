package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type condition struct {
	Field string      `json:"field"`
	Op    string      `json:"op"`
	Value interface{} `json:"value"`
}

/**
* ToJson
* @return et.Json
**/
func (s *condition) ToJson() et.Json {
	return et.Json{
		s.Field: et.Json{
			s.Op: s.Value,
		},
	}
}

/**
* Eq
* @param field string, value interface{}
* @return condition
**/
func Eq(field string, value interface{}) condition {
	return condition{
		Field: field,
		Op:    "eq",
		Value: value,
	}
}

/**
* Neg
* @param field string, value interface{}
* @return condition
**/
func Neg(field string, value interface{}) condition {
	return condition{
		Field: field,
		Op:    "ne",
		Value: value,
	}
}

/**
* Less
* @param field string, value interface{}
* @return condition
**/
func Less(field string, value interface{}) condition {
	return condition{
		Field: field,
		Op:    "less",
		Value: value,
	}
}

/**
* LessEq
* @param field string, value interface{}
* @return condition
**/
func LessEq(field string, value interface{}) condition {
	return condition{
		Field: field,
		Op:    "less_eq",
		Value: value,
	}
}

/**
* More
* @param field string, value interface{}
* @return condition
**/
func More(field string, value interface{}) condition {
	return condition{
		Field: field,
		Op:    "more",
		Value: value,
	}
}

/**
* MoreEq
* @param field string, value interface{}
* @return condition
**/
func MoreEq(field string, value interface{}) condition {
	return condition{
		Field: field,
		Op:    "more_eq",
		Value: value,
	}
}

/**
* Like
* @param field string, value interface{}
* @return condition
**/
func Like(field string, value interface{}) condition {
	return condition{
		Field: field,
		Op:    "like",
		Value: value,
	}
}

/**
* Ilike
* @param field string, value interface{}
* @return condition
**/
func Ilike(field string, value interface{}) condition {
	return condition{
		Field: field,
		Op:    "ilike",
		Value: value,
	}
}

/**
* In
* @param field string, value interface{}
* @return condition
**/
func In(field string, value interface{}) condition {
	return condition{
		Field: field,
		Op:    "in",
		Value: value,
	}
}

/**
* NotIn
* @param field string, value interface{}
* @return condition
**/
func NotIn(field string, value interface{}) condition {
	return condition{
		Field: field,
		Op:    "not_in",
		Value: value,
	}
}

/**
* Is
* @param field string, value interface{}
* @return condition
**/
func Is(field string, value interface{}) condition {
	return condition{
		Field: field,
		Op:    "is",
		Value: value,
	}
}

/**
* IsNot
* @param field string, value interface{}
* @return condition
**/
func IsNot(field string, value interface{}) condition {
	return condition{
		Field: field,
		Op:    "is_not",
		Value: value,
	}
}

/**
* Null
* @param field string
* @return condition
**/
func Null(field string) condition {
	return condition{
		Field: field,
		Op:    "null",
	}
}

/**
* NotNull
* @param field string
* @return condition
**/
func NotNull(field string) condition {
	return condition{
		Field: field,
		Op:    "not_null",
	}
}

/**
* Between
* @param field string, value []interface{}
* @return condition
**/
func Between(field string, value []interface{}) condition {
	return condition{
		Field: field,
		Op:    "between",
		Value: value,
	}
}

/**
* NotBetween
* @param field string, value []interface{}
* @return condition
**/
func NotBetween(field string, value []interface{}) condition {
	return condition{
		Field: field,
		Op:    "not_between",
		Value: value,
	}
}

type where struct {
	Wheres []et.Json `json:"where"`
}

func newWhere() *where {
	return &where{
		Wheres: []et.Json{},
	}
}

/**
* where
* @param cond condition
* @return *where
**/
func (s *where) where(cond condition, conector string) *where {
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
* @param cond condition
* @return *where
**/
func (s *where) Where(cond condition) *where {
	return s.where(cond, "and")
}

/**
* And
* @param cond condition
* @return *where
**/
func (s *where) And(cond condition) *where {
	return s.where(cond, "and")
}

/**
* Or
* @param cond condition
* @return *where
**/
func (s *where) Or(cond condition) *where {
	return s.where(cond, "or")
}
