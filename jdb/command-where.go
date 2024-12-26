package jdb

/**
* Where
* @param col interface{}
* @return *Command
**/
func (s *Command) Where(val interface{}) *LinqFilter {
	col := s.getColumn(val)
	if col != nil {
		return NewLinqFilter(s, col)
	}

	return NewLinqFilter(s, val)
}

/**
* And
* @param col interface{}
* @return *LinqWheres
**/
func (s *Command) And(col interface{}) *LinqFilter {
	result := s.Where(col)
	result.where.Conector = And
	return result
}

/**
* And
* @param col interface{}
* @return *LinqWhere
**/
func (s *Command) Or(col interface{}) *LinqFilter {
	result := s.Where(col)
	result.where.Conector = Or
	return result
}
