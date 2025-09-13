package jdb

import "github.com/cgalvisleon/et/et"

/**
* Where
* @param field string
* @return *Command
**/
func (s *Command) Where(val string) *Command {
	s.QlWhere.Where(val)

	return s
}

/**
* And
* @param val string
* @return *Command
**/
func (s *Command) And(val string) *Command {
	s.QlWhere.And(val)

	return s
}

/**
* And
* @param fval string
* @return *Command
**/
func (s *Command) Or(val string) *Command {
	s.QlWhere.Or(val)

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
* SetWhere
* @param tx *Tx
* @param params et.Json
* @return map[string]interface{}, error
**/
func (s *Command) SetWhere(wheres et.Json) *Command {
	s.QlWhere.SetWheres(wheres)

	return s
}
