package jdb

import "github.com/cgalvisleon/et/et"

/**
* Query
* @param query et.Json
* @return *Ql
**/
func (s *Model) Query(query et.Json) *Ql {
	result := s.db.From(s, "A")
	result.setQuery(query)
	return result
}

/**
* Select
* @param fields ...string
* @return *Ql
**/
func (s *Model) Select(fields ...string) *Ql {
	result := s.db.From(s, "A")
	result.IsDataSource = s.SourceField != ""
	return result.Select(fields...)
}

/**
* Data
* @param fields ...string
* @return *Ql
**/
func (s *Model) Data(fields ...string) *Ql {
	result := s.Select(fields...)
	result.IsDataSource = true
	return result
}

/**
* Where
* @param cond Condition
* @return *Ql
**/
func (s *Model) Where(cond *Condition) *Ql {
	result := s.db.From(s, "A")
	return result.Where(cond)
}

/**
* WhereByKeys
* @param data et.Json
* @return *Ql
**/
func (s *Model) WhereByKeys(data et.Json) *Ql {
	result := s.db.From(s, "A")
	for _, col := range s.PrimaryKeys {
		result.Where(Eq(col, data[col]))
	}
	return result
}

/**
* Join
* @param to *Model, as string, on Condition
* @return *Ql
**/
func (s *Model) Join(to *Model, as string, on *Condition) *Ql {
	result := s.db.From(s, "A")
	return result.Join(to, as, on)
}

/**
* Counted
* @return int, error
**/
func (s *Model) Counted() (int, error) {
	result := s.db.From(s, "A")
	return result.Counted()
}
