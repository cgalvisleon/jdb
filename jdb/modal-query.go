package jdb

import "github.com/cgalvisleon/et/et"

/**
* Query
* @param query et.Json
* @return *Ql
**/
func (s *Model) Query(query et.Json) *Ql {
	result := newQl(s.db)
	result.addFrom(s.Name)
	result.setQuery(query)

	return result
}
