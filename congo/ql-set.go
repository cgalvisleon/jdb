package jdb

import "github.com/cgalvisleon/et/et"

/**
* setQuery
* @param query et.Json
* @return *Ql
**/
func (s *Ql) setQuery(query et.Json) *Ql {
	s.setSelect(query.Json("select")).
		setWhere(query.Json("where")).
		setAnd(query.Json("and")).
		setOr(query.Json("or")).
		setOrder(query.Json("order")).
		setGroup(query.Json("group")).
		setHaving(query.Json("having")).
		setLimit(query.Json("limit"))

	return s
}

func (s *Ql) setSelect(selects et.Json) *Ql {
	s.Selects = selects
	return s
}

func (s *Ql) setWhere(where et.Json) *Ql {
	s.Wheres = where
	return s
}

func (s *Ql) setAnd(and et.Json) *Ql {
	s.Wheres = and
	return s
}

func (s *Ql) setOr(or et.Json) *Ql {
	s.Wheres = or
	return s
}

func (s *Ql) setJoin(join et.Json) *Ql {
	s.Joins = append(s.Joins, join)
	return s
}

func (s *Ql) setOrder(orderBy et.Json) *Ql {
	s.OrderBy = orderBy
	return s
}

func (s *Ql) setGroup(groupBy et.Json) *Ql {
	s.GroupBy = groupBy
	return s
}

func (s *Ql) setHaving(having et.Json) *Ql {
	s.Having = having
	return s
}

/**
* setLimit
* @param limits et.Json
* @return *Ql
**/
func (s *Ql) setLimit(limits et.Json) *Ql {
	s.Limit = limits
	return s
}
