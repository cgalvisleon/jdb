package jdb

/**
* getSelect
* @param name string
* @return *LinqSelect
**/
func (s *Command) getSelect(name string) *LinqSelect {
	field := NewField(name)
	if field == nil {
		return nil
	}

	from := &LinqFrom{
		Model: s.Model,
		As:    s.Model.Table,
	}

	return NewLinqSelect(from, field.Name)
}

/**
* Where
* @param field string
* @return *Command
**/
func (s *Command) Where(field string) *LinqFilter {
	col := s.getSelect(field)
	if col != nil {
		return NewLinqFilter(s, col)
	}

	return NewLinqFilter(s, field)
}

/**
* And
* @param field string
* @return *LinqWheres
**/
func (s *Command) And(field string) *LinqFilter {
	result := s.Where(field)
	result.where.Conector = And
	return result
}

/**
* And
* @param field string
* @return *LinqWhere
**/
func (s *Command) Or(field string) *LinqFilter {
	result := s.Where(field)
	result.where.Conector = Or
	return result
}
