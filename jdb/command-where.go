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
		field := s.GetField(v)
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

		if item["eq"] != nil {
			s.Eq(item["eq"])
		} else if item["neg"] != nil {
			s.Neg(item["neg"])
		} else if item["in"] != nil {
			s.In(item["in"])
		} else if item["like"] != nil {
			s.Like(item["like"])
		} else if item["more"] != nil {
			s.More(item["more"])
		} else if item["less"] != nil {
			s.Less(item["less"])
		} else if item["moreEq"] != nil {
			s.MoreEq(item["moreEq"])
		} else if item["lessEq"] != nil {
			s.LessEs(item["lessEq"])
		} else if item["search"] != nil {
			s.Search(item["search"])
		} else if item["between"] != nil {
			s.Between(item["between"])
		} else if item["isNull"] != nil {
			s.IsNull()
		} else if item["notNull"] != nil {
			s.NotNull()
		}
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
