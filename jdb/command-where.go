package jdb

type CommandFilter struct {
	Command *Command
	Wheres  []*LinqWhere
	where   *LinqWhere
}

/**
* NewCommandFilter
* @param l *Command
* @param wheres []*LinqWhere
* @return *CommandFilter
**/
func NewCommandFilter(c *Command, whers []*LinqWhere) *CommandFilter {
	return &CommandFilter{
		Command: c,
		Wheres:  whers,
		where:   &LinqWhere{},
	}
}

func (s *CommandFilter) add(operator Operator, val ...interface{}) *Command {
	s.where.Operator = operator
	col := s.Command.getColumn(val)
	if col != nil {
		s.where.B = col
	} else {
		s.where.B = val
	}
	s.Wheres = append(s.Wheres, s.where)

	return s.Command
}

/**
* Eq
* @param val interface{}
* @return *Command
**/
func (s *CommandFilter) Eq(val interface{}) *Command {
	return s.add(Equal, val)
}

/**
* Neg
* @param val interface{}
* @return *Command
**/
func (s *CommandFilter) Neg(val interface{}) *Command {
	return s.add(Neg, val)
}

/**
* In
* @param val ...interface{}
* @return *Command
**/
func (s *CommandFilter) In(val ...interface{}) *Command {
	return s.add(In, val)
}

/**
* Like
* @param val interface{}
* @return *Command
**/
func (s *CommandFilter) Like(val interface{}) *Command {
	return s.add(Like, val)
}

/**
* More
* @param val interface{}
* @return *Command
**/
func (s *CommandFilter) More(val interface{}) *Command {
	return s.add(More, val)
}

/**
* Less
* @param val interface{}
* @return *Command
**/
func (s *CommandFilter) Less(val interface{}) *Command {
	return s.add(Less, val)
}

/**
* MoreEq
* @param val interface{}
* @return *Command
**/
func (s *CommandFilter) MoreEq(val interface{}) *Command {
	return s.add(MoreEq, val)
}

/**
* LessEq
* @param val interface{}
* @return *Command
**/
func (s *CommandFilter) LessEs(val interface{}) *Command {
	return s.add(LessEq, val)
}

/**
* Between
* @param val1, val2 interface{}
* @return *Command
**/
func (s *CommandFilter) Between(val1, val2 interface{}) *Command {
	return s.add(Between, val1, val2)
}

/**
* IsNull
* @return *Command
**/
func (s *CommandFilter) IsNull() *Command {
	return s.add(IsNull, nil)
}

/**
* Select
* @param columns ...interface{}
* @return *Command
**/
func (s *Command) Where(col interface{}) *CommandFilter {
	where := &LinqWhere{
		Conector: Not,
	}

	_col := s.getColumn(col)
	if _col != nil {
		where.A = _col
	} else {
		where.A = col
	}

	result := &CommandFilter{
		Command: s,
		Wheres:  s.Wheres,
		where:   where,
	}

	return result
}

/**
* And
* @param col interface{}
* @return *LinqWheres
**/
func (s *Command) And(col interface{}) *CommandFilter {
	result := s.Where(col)
	result.where.Conector = And
	return result
}

/**
* And
* @param col interface{}
* @return *LinqWhere
**/
func (s *Command) Or(col interface{}) *CommandFilter {
	result := s.Where(col)
	result.where.Conector = Or
	return result
}
