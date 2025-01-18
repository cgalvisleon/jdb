package jdb

import (
	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
)

/**
* calcOffset
* @return *Ql
**/
func (s *Ql) calcOffset() *Ql {
	s.Offset = (s.Sheet - 1) * s.Limit
	if s.Offset < 0 {
		s.Offset = 0
	}

	return s
}

/**
* SetLimit
* @param limit int
* @return *Ql
**/
func (s *Ql) setLimit(limit int) *Ql {
	max := envar.GetInt(1000, "QUERY_LIMIT")
	if limit > max {
		limit = max
	}
	s.Limit = limit

	return s
}

/**
* SetPage
* @param page int
* @return *Ql
**/
func (s *Ql) setPage(page int) *Ql {
	return s.Page(page)
}

/**
* listLimit
* @return interface{}
**/
func (s *Ql) listLimit() interface{} {
	if s.Sheet > 0 {
		return et.Json{
			"limit": s.Limit,
			"page":  s.Sheet,
		}
	}

	return s.Limit
}
