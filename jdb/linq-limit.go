package jdb

import "github.com/cgalvisleon/et/et"

/**
* calcOffset
* @return *Linq
**/
func (s *Linq) calcOffset() *Linq {
	s.Offset = (s.Sheet - 1) * s.Limit
	if s.Offset < 0 {
		s.Offset = 0
	}

	return s
}

/**
* SetLimit
* @param limit int
* @return *Linq
**/
func (s *Linq) setLimit(limit int) *Linq {
	s.Limit = limit
	return s
}

/**
* SetPage
* @param page int
* @return *Linq
**/
func (s *Linq) setPage(page int) *Linq {
	return s.Page(page)
}

/**
* listLimit
* @return interface{}
**/
func (s *Linq) listLimit() interface{} {
	if s.Sheet > 0 {
		return et.Json{
			"limit": s.Limit,
			"page":  s.Sheet,
		}
	}

	return s.Limit
}
