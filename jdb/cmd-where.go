package jdb

import "github.com/cgalvisleon/et/et"

/**
* Where
* @param cond Condition
* @return *Cmd
**/
func (s *Cmd) Where(cond *Condition) *Cmd {
	s.where.Where(cond)
	return s
}

/**
* And
* @param cond Condition
* @return *Cmd
**/
func (s *Cmd) And(cond *Condition) *Cmd {
	s.where.And(cond)
	return s
}

/**
* Or
* @param cond Condition
* @return *Cmd
**/
func (s *Cmd) Or(cond *Condition) *Cmd {
	s.where.Or(cond)
	return s
}

/**
* WhereByPrimaryKeys
* @param data et.Json
* @return *Cmd
**/
func (s *Cmd) WhereByPrimaryKeys(data et.Json) *Cmd {
	model := s.Froms[0].Model
	for _, col := range model.PrimaryKeys {
		s.Where(Eq(col, data[col]))
	}
	return s
}
