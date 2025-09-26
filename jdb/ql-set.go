package jdb

import "github.com/cgalvisleon/et/et"

/**
* setQuery
* @param query et.Json
* @return *Ql
**/
func (s *Ql) setQuery(query et.Json) *Ql {
	s.setType(query.String("type")).
		setDebug(query.Bool("debug")).
		setSelect(query.Json("select")).
		setAtribs(query.Json("atribs")).
		setRollup(query.Json("rollups")).
		setRelation(query.Json("relations")).
		setCall(query.Json("calls")).
		setJoin(query.ArrayJson("joins")).
		setWhere(query.ArrayJson("where")).
		setOrderBy(query.Json("order_by")).
		setGroupBy(query.ArrayStr("group_by")).
		setHaving(query.ArrayJson("having")).
		setLimit(query.Json("limit"))

	return s
}

/**
* setType
* @param tp string
* @return *Ql
**/
func (s *Ql) setType(tp string) *Ql {
	if _, ok := TpQuerys[tp]; !ok {
		return s
	}

	s.Type = tp
	return s
}

/**
* setDebug
* @param debug bool
* @return *Ql
**/
func (s *Ql) setDebug(debug bool) *Ql {
	s.isDebug = debug
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
* setAtribs
* @param atribs et.Json
* @return *Ql
**/
func (s *Ql) setAtribs(atribs et.Json) *Ql {
	s.Atribs = atribs
	return s
}

/**
* setRollup
* @param rollups et.Json
* @return *Ql
**/
func (s *Ql) setRollup(rollups et.Json) *Ql {
	s.Rollups = rollups
	return s
}

/**
* setRelation
* @param relations et.Json
* @return *Ql
**/
func (s *Ql) setRelation(relations et.Json) *Ql {
	s.Relations = relations
	return s
}

/**
* setCall
* @param calls et.Json
* @return *Ql
**/
func (s *Ql) setCall(calls et.Json) *Ql {
	s.Calls = calls
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
* setOrderBy
* @param orderBy et.Json
* @return *Ql
**/
func (s *Ql) setOrderBy(orderBy et.Json) *Ql {
	s.OrderBy = orderBy
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
* @return *Ql
**/
func (s *Ql) setLimit(limits et.Json) *Ql {
	s.Limit = limits
	return s
}
