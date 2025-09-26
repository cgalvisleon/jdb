package jdb

import "github.com/cgalvisleon/et/et"

/**
* Query
* @param query et.Json
* @return *Ql
**/
func (s *Model) Query(query et.Json) *Ql {
	result := newQl(s.db)
	result.addFrom(s.Name, "A")
	return result.setQuery(query)
}

/**
* Select
* @param fields interface{}
* @return *Ql
**/
func (s *Model) Select(fields interface{}) *Ql {
	result := newQl(s.db)
	result.addFrom(s.Name, "A")
	return result.Select(fields)
}

/**
* Object
* @param fields interface{}
* @return *Ql
**/
func (s *Model) Object(fields interface{}) *Ql {
	result := newQl(s.db)
	result.addFrom(s.Name, "A")
	return result.Object(fields)
}

/**
* Where
* @param cond Condition
* @return *Ql
**/
func (s *Model) Where(cond Condition) *Ql {
	result := newQl(s.db)
	result.addFrom(s.Name, "A")
	if s.SourceField != "" {
		result.Type = TpObject
	}

	return result.Where(cond)
}
