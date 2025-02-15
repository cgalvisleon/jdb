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
		field := s.From.GetField(v)
		if field != nil {
			s.where = NewQlWhere(field)
			return s
		}
	}

	s.where = NewQlWhere(val)
	return s
}

/**
* And
* @param val interface{}
* @return *QlFilter
**/
func (s *Command) And(val interface{}) *QlFilter {
	result := s.Where(val)
	result.where.Conector = And
	return s.QlFilter
}

/**
* And
* @param field string
* @return *QlFilter
**/
func (s *Command) Or(val interface{}) *QlFilter {
	result := s.Where(val)
	result.where.Conector = Or
	return s.QlFilter
}

/**
* Select
* @param fields ...string
* @return *Ql
**/
func (s *Command) Select(fields ...string) *Ql {
	return nil
}

/**
* Data
* @param fields ...string
* @return *Ql
**/
func (s *Command) Data(fields ...string) *Ql {
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
