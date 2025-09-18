package jdb

import (
	"github.com/cgalvisleon/et/et"
)

type JQuery struct {
	selects et.Json
	froms   []map[string]string
	wheres  et.Json
	joins   et.Json
	orders  et.Json
	groups  et.Json
	having  et.Json
	limits  et.Json
}

func NewJQuery() *JQuery {
	return &JQuery{
		froms:   []map[string]string{},
		selects: et.Json{},
		wheres:  et.Json{},
		joins:   et.Json{},
		orders:  et.Json{},
		groups:  et.Json{},
		having:  et.Json{},
		limits:  et.Json{},
	}
}

func (s *JQuery) ToJson() et.Json {
	return et.Json{
		"froms":   s.froms,
		"selects": s.selects,
		"wheres":  s.wheres,
		"joins":   s.joins,
		"orders":  s.orders,
		"groups":  s.groups,
		"limits":  s.limits,
	}
}

func (s *JQuery) setFrom(name string) *JQuery {
	n := len(s.froms)
	as := string(rune(65 + n))
	s.froms = append(s.froms, map[string]string{
		name: as,
	})

	if n != 0 {
		s.joins = et.Json{
			"from": name,
		}
	}
	return s
}

func (s *JQuery) setSelect(selects et.Json) *JQuery {
	s.selects = selects
	return s
}

func (s *JQuery) setWhere(where et.Json) *JQuery {
	s.wheres = where
	return s
}

func (s *JQuery) setAnd(and et.Json) *JQuery {
	s.wheres = and
	return s
}

func (s *JQuery) setOr(or et.Json) *JQuery {
	s.wheres = or
	return s
}

func (s *JQuery) setJoin(joins et.Json) *JQuery {
	s.joins = joins
	return s
}

func (s *JQuery) setOrder(orders et.Json) *JQuery {
	s.orders = orders
	return s
}

func (s *JQuery) setGroup(groups et.Json) *JQuery {
	s.groups = groups
	return s
}

func (s *JQuery) setHaving(having et.Json) *JQuery {
	s.having = having
	return s
}

func (s *JQuery) setLimit(limits et.Json) *JQuery {
	s.limits = limits
	return s
}

func (s *JQuery) sql() string {
	return s.ToJson().String()
}

func Query(query et.Json) (string, error) {
	j := NewJQuery()

	for k := range query {
		jq := query.Json(k)
		j.setFrom(k).
			setSelect(jq.Json("select")).
			setWhere(jq.Json("where")).
			setAnd(jq.Json("and")).
			setOr(jq.Json("or")).
			setOrder(jq.Json("order")).
			setGroup(jq.Json("group")).
			setHaving(jq.Json("having")).
			setLimit(jq.Json("limit"))
	}

	return j.sql(), nil
}
