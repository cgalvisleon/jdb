package jdb

import "github.com/cgalvisleon/et/et"

func (s *Linq) All() (et.Items, error) {
	return (*s.Db.Driver).Query(s)
}

/**
* Offset
* @param offset int
* @return *Linq
**/
func (s *Linq) Page(val int) *Linq {
	s.page = val
	s.Offset = s.Limit * (s.page - 1)
	return s
}

/**
* Limit
* @param limit int
* @return *Linq
**/
func (s *Linq) Rows(val int) *Linq {
	s.Limit = val
	s.Offset = s.Limit * (s.page - 1)
	return s
}
