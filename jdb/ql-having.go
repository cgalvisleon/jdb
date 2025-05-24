package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/et"
)

type QlHaving struct {
	*QlWhere
	Ql *Ql
}

/**
* NewQlHaving
* @param ql *Ql
* @return *QlHaving
**/
func NewQlHaving(ql *Ql) *QlHaving {
	return &QlHaving{
		Ql:      ql,
		QlWhere: newQlWhere(ql.validator),
	}
}

/**
* Having
* @param val string
* @return *QlHaving
**/
func (s *QlHaving) And(val string) *QlHaving {
	field := s.Ql.getColumnField(val)
	if field != nil {
		s.setAnd(val)
	}

	return s
}

/**
* Or
* @param val string
* @return *QlHaving
**/
func (s *QlHaving) Or(val string) *QlHaving {
	field := s.Ql.getColumnField(val)
	if field != nil {
		s.setOr(val)
	}

	return s
}

/**
* Select
* @param fields ...interface{}
* @return *Ql
**/
func (s *QlHaving) Select(fields ...interface{}) *Ql {
	return s.Ql.Select(fields...)
}

/**
* Data
* @param fields ...interface{}
* @return *Ql
**/
func (s *QlHaving) Data(fields ...interface{}) *Ql {
	return s.Ql.Data(fields...)
}

/**
* Having
* @param field string
* @return *QlWhere
**/
func (s *Ql) Having(val string) *QlHaving {
	field := s.getColumnField(val)
	if field != nil {
		s.Havings.Where(field)
	}

	return s.Havings
}

/**
* setHavings
* @param havings et.Json
* @return *Ql
**/
func (s *Ql) setHavings(havings et.Json) *Ql {
	if len(havings) == 0 {
		return s
	}

	and := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				s.Havings.setAnd(key)
				s.Havings.setValue(val.Json(key))
			}
		}
	}

	or := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				s.Havings.setOr(key)
				s.Havings.setValue(val.Json(key))
			}
		}
	}

	for key := range havings {
		if slices.Contains([]string{"and", "AND", "or", "OR"}, key) {
			continue
		}

		s.Having(key).setValue(havings.Json(key))
	}

	for key := range havings {
		switch key {
		case "and", "AND":
			vals := havings.ArrayJson(key)
			and(vals)
		case "or", "OR":
			vals := havings.ArrayJson(key)
			or(vals)
		}
	}

	return s
}

/**
* getHavings
* @return et.Json
**/
func (s *Ql) getHavings() et.Json {
	return s.Havings.getWheres()
}
