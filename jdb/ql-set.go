package jdb

import "github.com/cgalvisleon/et/et"

/**
* setQuery
* @param query et.Json
* @return *Ql
**/
func (s *Ql) setQuery(query et.Json) *Ql {
	s.SetDebug(query.Bool("debug")).
		setSelect(query.Json("select")).
		setJoin(query.ArrayJson("joins")).
		setWhere(query.ArrayJson("where")).
		setHiddens(query.ArrayStr("hidden")).
		setOrderBy(query.Json("order_by")).
		setGroupBy(query.ArrayStr("group_by")).
		setHaving(query.ArrayJson("having")).
		setLimit(query.Json("limit"))
	return s
}

/**
* setDebug
* @param debug bool
* @return *Ql
**/
func (s *Ql) SetDebug(debug bool) *Ql {
	s.IsDebug = debug
	return s
}

/**
* setSelect
* @param selects et.Json
* @return *Ql
**/
func (s *Ql) setSelect(selects et.Json) *Ql {
	s.Selects = selects
	return s
}

/**
* setJoin
* @param joins []et.Json
* @return *Ql
**/
func (s *Ql) setJoin(joins []et.Json) *Ql {
	s.Joins = joins
	return s
}

/**
* setWhere
* @param where []et.Json
* @return *Ql
**/
func (s *Ql) setWhere(where []et.Json) *Ql {
	s.Wheres = where
	return s
}

/**
* setHiddens
* @param hiddens []string
* @return *Ql
**/
func (s *Ql) setHiddens(hiddens []string) *Ql {
	s.Hiddens = hiddens
	return s
}

/**
* setOrderBy
* @param orderBy et.Json
* @return *Ql
**/
func (s *Ql) setOrderBy(orderBy et.Json) *Ql {
	s.OrdersBy = orderBy
	return s
}

/**
* setGroupBy
* @param groupBy []string
* @return *Ql
**/
func (s *Ql) setGroupBy(groupBy []string) *Ql {
	s.GroupBy = groupBy
	return s
}

/**
* setHaving
* @param having []et.Json
* @return *Ql
**/
func (s *Ql) setHaving(having []et.Json) *Ql {
	s.Havings = having
	return s
}

/**
* setLimit
* @param limits et.Json
**/
func (s *Ql) setLimit(limits et.Json) *Ql {
	s.Limits = limits
	return s
}

/**
* Result executes the query and returns the result
* @return et.Items, error
**/
func (s *Ql) Result() (et.Items, error) {
	page := s.Limits.ValInt(1, "page")
	rows := s.Limits.ValInt(1000, "rows")
	return s.Limit(page, rows)
}
