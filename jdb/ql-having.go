package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type QlHaving struct {
	*QlFilter
	Ql *Ql
}

/**
* Having
* @param field string
* @return *QlWhere
**/
func (s *Ql) Having(field string) *QlHaving {
	return s.Havings.on(field)
}

/**
* Having
* @param field string
* @return *QlWhere
**/
func (s *QlHaving) on(field string) *QlHaving {
	sel := s.Ql.GetSelect(field)
	if sel != nil {
		s.where = NewQlWhere(sel)
	} else {
		s.where = NewQlWhere(field)
	}

	return s
}

/**
* And
* @param field string
* @return *QlFilter
**/
func (s *QlHaving) And(val interface{}) *QlFilter {
	field, ok := val.(string)
	if ok {
		result := s.on(field)
		result.where.Conector = And
	}

	return s.QlFilter
}

/**
* Or
* @param field string
* @return *QlFilter
**/
func (s *QlHaving) Or(val interface{}) *QlFilter {
	field, ok := val.(string)
	if ok {
		result := s.on(field)
		result.where.Conector = Or
	}

	return s.QlFilter
}

/**
* Select
* @param fields ...string
* @return *Ql
**/
func (s *QlHaving) Select(fields ...string) *Ql {
	return s.Ql
}

/**
* Data
* @param fields ...string
* @return *Ql
**/
func (s *QlHaving) Data(fields ...string) *Ql {
	return s.Ql
}

/**
* Exec
* @return et.Items, error
**/
func (s *QlHaving) Exec() (et.Items, error) {
	return et.Items{}, nil
}

/**
* One
* @return et.Item, error
**/
func (s *QlHaving) One() (et.Item, error) {
	return et.Item{}, nil
}

/**
* setHavings
* @param havings []et.Json
* @return *Ql
**/
func (s *Ql) setHavings(havings []et.Json) *Ql {
	for _, val := range havings {
		from := val.Str("from")
		model := models[from]
		if model != nil {
			on := val.Json("on")
			key := strs.Format(`%s.%s`, from, on.Str("key"))
			to := on.Str("to")
			foreign := on.Str("foreignKey")
			foreignKey := strs.Format(`%s.%s`, to, foreign)
			s.Join(model).On(key).
				Eq(foreignKey)
		}
	}

	return s
}

/**
* listHavings
* @return []string
**/
func (s *Ql) listHavings() []string {
	result := []string{}
	for _, val := range s.Havings.Wheres {
		result = append(result, val.String())
	}

	return result
}
