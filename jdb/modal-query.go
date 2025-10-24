package jdb

import "github.com/cgalvisleon/et/et"

/**
* Query
* @param query et.Json
* @return *Ql
**/
func (s *Model) Query(query et.Json) *Ql {
	result := s.db.From(s)
	return result.setQuery(query)
}

/**
* Select
* @param fields ...string
* @return *Ql
**/
func (s *Model) Select(fields ...string) *Ql {
	result := s.db.From(s)
	return result.Select(fields...)
}

/**
* Where
* @param cond Condition
* @return *Ql
**/
func (s *Model) Where(cond *Condition) *Ql {
	result := s.db.From(s)
	return result.Where(cond)
}

/**
* Join
* @param to *Model, as string, on Condition
* @return *Ql
**/
func (s *Model) Join(to *Model, as string, on *Condition) *Ql {
	result := s.db.From(s)
	return result.Join(to, as, on)
}

/**
* Counted
* @return int, error
**/
func (s *Model) Counted() (int, error) {
	result := s.db.From(s)
	return result.Counted()
}
