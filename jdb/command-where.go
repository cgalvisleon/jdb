package jdb

func (s *Command) getField(name string) *Field {
	return s.From.getField(name)
}

/**
* Where
* @param field string
* @return *Command
**/
func (s *Command) Where(val string) *Command {
	field := s.getField(val)
	if field != nil {
		s.where(field)
	}

	return s
}

/**
* And
* @param val string
* @return *Command
**/
func (s *Command) And(val string) *Command {
	field := s.getField(val)
	if field != nil {
		s.and(field)
	}

	return s
}

/**
* And
* @param fval string
* @return *Command
**/
func (s *Command) Or(val string) *Command {
	field := s.getField(val)
	if field != nil {
		s.or(field)
	}

	return s
}

/**
* Eq
* @param val interface{}
* @return *Command
**/
func (s *Command) Eq(val interface{}) *Command {
	s.QlWhere.Eq(val)

	return s
}

/**
* Neg
* @param val interface{}
* @return *Command
**/
func (s *Command) Neg(val interface{}) *Command {
	s.QlWhere.Neg(val)

	return s
}

/**
* In
* @param val ...any
* @return *Command
**/
func (s *Command) In(val ...any) *Command {
	s.QlWhere.In(val...)

	return s
}

/**
* Like
* @param val interface{}
* @return *Command
**/
func (s *Command) Like(val interface{}) *Command {
	s.QlWhere.Like(val)

	return s
}

/**
* More
* @param val interface{}
* @return *Command
**/
func (s *Command) More(val interface{}) *Command {
	s.QlWhere.More(val)

	return s
}

/**
* Less
* @param val interface{}
* @return *Command
**/
func (s *Command) Less(val interface{}) *Command {
	s.QlWhere.Less(val)

	return s
}

/**
* MoreEq
* @param val interface{}
* @return *Command
**/
func (s *Command) MoreEq(val interface{}) *Command {
	s.QlWhere.MoreEq(val)

	return s
}

/**
* LessEq
* @param val interface{}
* @return *Command
**/
func (s *Command) LessEq(val interface{}) *Command {
	s.QlWhere.LessEq(val)

	return s
}

/*
*
* Between
* @param vals interface{}
* @return *Command
**/
func (s *Command) Between(vals interface{}) *Command {
	s.QlWhere.Between(vals)

	return s
}

/**
* Search
* @param language string, val interface{}
* @return *Command
**/
func (s *Command) Search(language string, val interface{}) *Command {
	s.QlWhere.Search(language, val)

	return s
}

/**
* IsNull
* @return *Command
**/
func (s *Command) IsNull() *Command {
	s.QlWhere.IsNull()

	return s
}

/**
* NotNull
* @return *Command
**/
func (s *Command) NotNull() *Command {
	s.QlWhere.NotNull()

	return s
}

/**
* History
* @param v bool
* @return *Command
**/
func (s *Command) History(v bool) *Command {
	s.QlWhere.History(v)

	return s
}

/**
* Debug
* @param v bool
* @return *Command
**/
func (s *Command) Debug() *Command {
	s.QlWhere.Debug()

	return s
}
