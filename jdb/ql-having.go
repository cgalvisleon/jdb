package jdb

import (
	"github.com/cgalvisleon/et/et"
)

type QlHaving struct {
	*QlWhere
	Ql *Ql
}

/**
* Having
* @param val interface{}
* @return *QlHaving
**/
func (s *QlHaving) And(val interface{}) *QlHaving {
	switch v := val.(type) {
	case string:
		field := s.Ql.getField(v)
		if field != nil {
			s.and(field)
			return s
		}
	}

	return s
}

/**
* Or
* @param val interface{}
* @return *QlHaving
**/
func (s *QlHaving) Or(val interface{}) *QlHaving {
	switch v := val.(type) {
	case string:
		field := s.Ql.getField(v)
		if field != nil {
			s.or(field)
			return s
		}
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
	}

	return s.Havings
}

/**
* setHavings
* @param havings et.Json
* @return *Ql
**/
func (s *Ql) setHavings(havings et.Json) *Ql {
	for key := range havings {
		val := havings.Json(key)
		s.Having(key).
			setValue(val)
	}

	return s
}

/**
* listHavings
* @return et.Json
**/
func (s *Ql) listHavings() et.Json {
	return s.Havings.listWheres(s.asField)
}
