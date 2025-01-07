package jdb

/**
* Where
* @param field string
* @return *Command
**/
func (s *Command) Where(val interface{}) *Command {
	switch v := val.(type) {
	case string:
		field := s.GetField(v)
		if field != nil {
			s.where = NewLinqWhere(field)
			return s
		}
	}

	s.where = NewLinqWhere(val)
	return s
}

/**
* And
* @param val interface{}
* @return *LinqFilter
**/
func (s *Command) And(val interface{}) *LinqFilter {
	result := s.Where(val)
	result.where.Conector = And
	return s.LinqFilter
}

/**
* And
* @param field string
* @return *LinqFilter
**/
func (s *Command) Or(val interface{}) *LinqFilter {
	result := s.Where(val)
	result.where.Conector = Or
	return s.LinqFilter
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
		sel := s.GetReturn(name)
		if sel != nil {
			s.Returns = append(s.Returns, sel)
		}
	}

	return s
}
