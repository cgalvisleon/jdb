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
* SetPage
* @param page int
* @return *Ql
**/
func (s *Ql) setPage(page int) *Ql {
	if page > 0 {
		s.Page(page)
	}

	return s
}

/**
* SetLimit
* @param limit int
* @return *Ql
**/
func (s *Ql) setLimit(limit int) (interface{}, error) {
	max := envar.GetInt(1000, "QUERY_LIMIT")
	if limit > max {
		limit = max
	}
	s.Limit = limit
	if s.Limit == 0 {
		return s.All()
	} else if s.Limit == 1 {
		return s.One()
	} else if s.Sheet > 0 {
		return s.Rows(s.Limit)
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
