package jdb

import "github.com/cgalvisleon/et/et"

/**
* Where
* @param field string
* @return *Command
**/
func (s *Command) Where(val interface{}) *Command {
	switch v := val.(type) {
	case string:
		field := s.From.GetField(v, false)
		if field != nil {
			s.where = NewLinqWhere(field)
			return s
		}
	}

	s.where = NewLinqWhere(val)
	return s
}

/**
* And
* @param val interface{}
* @return *LinqFilter
**/
func (s *Command) And(val interface{}) *LinqFilter {
	result := s.Where(val)
	result.where.Conector = And
	return s.LinqFilter
}

/**
* And
* @param field string
* @return *LinqFilter
**/
func (s *Command) Or(val interface{}) *LinqFilter {
	result := s.Where(val)
	result.where.Conector = Or
	return s.LinqFilter
}

/**
* Select
* @param fields ...string
* @return *Linq
**/
func (s *Command) Select(fields ...string) *Linq {
	return nil
}

/**
* Data
* @param fields ...string
* @return *Linq
**/
func (s *Command) Data(fields ...string) *Linq {
	return nil
}

func (s *Command) setWhere(wheres []et.Json) *Command {
	for _, item := range wheres {
		if item["key"] != nil {
			s.Where(item["key"])
		} else if item["and"] != nil {
			s.And(item["and"])
		} else if item["or"] != nil {
			s.Or(item["or"])
		}

		s.setCondition(item)
	}

	return s
}

/**
* listWheres
* @return []string
**/
func (s *Command) listWheres() []string {
	result := []string{}
	for _, val := range s.Wheres {
		result = append(result, val.String())
	}

	return result
}
