package jdb

import (
	"slices"
	"strings"

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
* SetHavings
* @param havings et.Json
* @return *Ql
**/
func (s *Ql) SetHavings(havings et.Json) *Ql {
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
		key = strings.ToLower(key)
		if slices.Contains([]string{"and", "or"}, key) {
			continue
		}

		val := havings.Json(key)
		s.Having(key).setValue(val)
	}

	for key := range havings {
		switch strings.ToLower(key) {
		case "and":
			vals := havings.ArrayJson(key)
			and(vals)
		case "or":
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
