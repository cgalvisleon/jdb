package jdb

import "github.com/cgalvisleon/et/et"

/**
* setQuery
* @param query et.Json
* @return *Ql
**/
func (s *Ql) setQuery(query et.Json) *Ql {
	s.setType(query.String("type")).
		setSelect(query.Json("select")).
		setAtribs(query.Json("atribs")).
		setRollup(query.Json("rollups")).
		setRelation(query.Json("relations")).
		setCall(query.Json("calls")).
		setJoin(query.ArrayJson("joins")).
		setWhere(query.Json("where")).
		setAnd(query.Json("and")).
		setOr(query.Json("or")).
		setOrder(query.Json("order")).
		setGroup(query.Json("group")).
		setHaving(query.Json("having")).
		setLimit(query.Json("limit"))

	return s
}

func (s *Ql) setType(tp string) *Ql {
	if _, ok := TpQuerys[tp]; !ok {
		return s
	}

	s.Type = tp
	return s
}

func (s *Ql) setSelect(selects et.Json) *Ql {
	s.Selects = selects
	return s
}

func (s *Ql) setAtribs(atribs et.Json) *Ql {
	s.Atribs = atribs
	return s
}

func (s *Ql) setRollup(rollups et.Json) *Ql {
	s.Rollups = rollups
	return s
}

func (s *Ql) setRelation(relations et.Json) *Ql {
	s.Relations = relations
	return s
}

func (s *Ql) setCall(calls et.Json) *Ql {
	s.Calls = calls
	return s
}

func (s *Ql) setJoin(joins []et.Json) *Ql {
	s.Joins = joins
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
