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
* Having
* @param val string
* @return *QlHaving
**/
func (s *QlHaving) And(val string) *QlHaving {
	field := s.Ql.getField(val)
	if field != nil {
		s.and(field)
	} else {
		s.and(val)
	}

	return s
}

/**
* Or
* @param val string
* @return *QlHaving
**/
func (s *QlHaving) Or(val string) *QlHaving {
	field := s.Ql.getField(val)
	if field != nil {
		s.or(field)
	} else {
		s.or(val)
	}

	return s
}

/**
* Select
* @param fields ...string
* @return *Ql
**/
func (s *QlHaving) Select(fields ...string) *Ql {
	return s.Ql.Select(fields...)
}

/**
* Data
* @param fields ...string
* @return *Ql
**/
func (s *QlHaving) Data(fields ...string) *Ql {
	return s.Ql.Data(fields...)
}

/**
* Having
* @param field string
* @return *QlWhere
**/
func (s *Ql) Having(val string) *QlHaving {
	field := s.getField(val)
	if field != nil {
		s.Havings.where(field)
	} else {
		s.Havings.where(val)
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
				s.Havings.and(key)
				s.Havings.setValue(val.Json(key), s.validator)
			}
		}
	}

	or := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				s.Havings.or(key)
				s.Havings.setValue(val.Json(key), s.validator)
			}
		}
	}

	for key := range havings {
		if slices.Contains([]string{"and", "AND", "or", "OR"}, key) {
			continue
		}

		s.Having(key).setValue(havings.Json(key), s.validator)
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
* listHavings
* @return et.Json
**/
func (s *Ql) listHavings() et.Json {
	return s.Havings.listWheres()
}
