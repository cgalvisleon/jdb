package jdb

import "github.com/cgalvisleon/et/et"

/**
* setLimit
* @param limits et.Json
* @return *JQuery
**/
func (s *JQuery) setLimit(limits et.Json) *JQuery {
	s.Limit = limits
	return s
}

/**
* setQuery
* @param query et.Json
* @return *JQuery
**/
func (s *JQuery) setQuery(query et.Json) *JQuery {
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

func (s *JQuery) setSelect(selects et.Json) *JQuery {
	s.Selects = selects
	return s
}

func (s *JQuery) setWhere(where et.Json) *JQuery {
	s.Wheres = where
	return s
}

func (s *JQuery) setAnd(and et.Json) *JQuery {
	s.Wheres = and
	return s
}

func (s *JQuery) setOr(or et.Json) *JQuery {
	s.Wheres = or
	return s
}

func (s *JQuery) setJoin(join et.Json) *JQuery {
	s.Joins = append(s.Joins, join)
	return s
}

func (s *JQuery) setOrder(orderBy et.Json) *JQuery {
	s.OrderBy = orderBy
	return s
}

func (s *JQuery) setGroup(groupBy et.Json) *JQuery {
	s.GroupBy = groupBy
	return s
}

func (s *JQuery) setHaving(having et.Json) *JQuery {
	s.Having = having
	return s
}
