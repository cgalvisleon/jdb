package jdb

/**
* GetSelect
* @param name string
* @return *LinqSelect
**/
func (s *Command) GetSelect(name string) *LinqSelect {
	field := s.GetField(name)
	if field == nil {
		return nil
	}

	from := &LinqFrom{
		Model: s.Model,
		As:    s.Model.Table,
	}

	return NewLinqSelect(from, field)
}

/**
* Where
* @param field string
* @return *LinqFilter
**/
func (s *Command) Where(field string) *LinqFilter {
	col := s.GetSelect(field)
	if col != nil {
		return NewLinqFilter(s, col)
	}

	return NewLinqFilter(s, field)
}

/**
* And
* @param val interface{}
* @return *LinqFilter
**/
func (s *Command) And(val interface{}) *LinqFilter {
	field, ok := val.(string)
	if !ok {
		return nil
	}

	result := s.Where(field)
	result.where.Conector = And
	return result
}

/**
* And
* @param field string
* @return *LinqFilter
**/
func (s *Command) Or(val interface{}) *LinqFilter {
	field, ok := val.(string)
	if !ok {
		return nil
	}

	result := s.Where(field)
	result.where.Conector = Or
	return result
}

/**
* Select
* @param fields ...string
* @return *Linq
**/
func (s *Command) Select(fields ...string) *Linq {
	return nil
}

/**
* Data
* @param fields ...string
* @return *Linq
**/
func (s *Command) Data(fields ...string) *Linq {
	return nil
}

/**
* Return
* @param fields ...string
* @return *Command
**/
func (s *Command) Return(fields ...string) *Command {
	for _, name := range fields {
		sel := s.GetSelect(name)
		if sel != nil {
			s.Returns = append(s.Returns, sel)
		}
	}

	return s
}
