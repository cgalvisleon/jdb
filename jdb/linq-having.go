package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type LinqHaving struct {
	*LinqFilter
	Linq *Linq
}

/**
* Having
* @param field string
* @return *LinqWhere
**/
func (s *Linq) Having(field string) *LinqHaving {
	return s.Havings.on(field)
}

/**
* Having
* @param field string
* @return *LinqWhere
**/
func (s *LinqHaving) on(field string) *LinqHaving {
	sel := s.Linq.GetSelect(field)
	if sel != nil {
		s.where = NewLinqWhere(sel)
	} else {
		s.where = NewLinqWhere(field)
	}

	return s
}

/**
* And
* @param field string
* @return *LinqFilter
**/
func (s *LinqHaving) And(val interface{}) *LinqFilter {
	field, ok := val.(string)
	if ok {
		result := s.on(field)
		result.where.Conector = And
	}

	return s.LinqFilter
}

/**
* Or
* @param field string
* @return *LinqFilter
**/
func (s *LinqHaving) Or(val interface{}) *LinqFilter {
	field, ok := val.(string)
	if ok {
		result := s.on(field)
		result.where.Conector = Or
	}

	return s.LinqFilter
}

/**
* Select
* @param fields ...string
* @return *Linq
**/
func (s *LinqHaving) Select(fields ...string) *Linq {
	return s.Linq
}

/**
* Data
* @param fields ...string
* @return *Linq
**/
func (s *LinqHaving) Data(fields ...string) *Linq {
	return s.Linq
}

/**
* Return
* @param fields ...string
* @return *Command
**/
func (s *LinqHaving) Return(fields ...string) *Command {
	return &Command{}
}

/**
* setHavings
* @param havings []et.Json
* @return *Linq
**/
func (s *Linq) setHavings(havings []et.Json) *Linq {
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
func (s *Linq) listHavings() []string {
	result := []string{}
	for _, val := range s.Havings.Wheres {
		result = append(result, val.String())
	}

	return result
}
