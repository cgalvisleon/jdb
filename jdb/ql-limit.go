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
	max := envar.GetInt(1000, "QUERY_LIMIT")
	if s.Limit > max {
		s.Limit = max
	}

	s.Offset = (s.Sheet - 1) * s.Limit
	if s.Offset < 0 {
		s.Offset = 0
	}

	return s
}

/**
* SetPage
* @param page int
* @return *Ql
**/
func (s *Ql) setPage(page int) *Ql {
	s.Page(page)

	return s
}

/**
* SetLimit
* @param limit int
* @return *Ql
**/
func (s *Ql) setLimit(limit int) (interface{}, error) {
	s.Limit = limit
	if s.Limit <= 0 {
		return s.All()
	} else if s.Limit == 1 {
		return s.One()
	}

	return s.First(s.Limit)
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
