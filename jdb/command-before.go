package jdb

/**
* BeforeInsert
* @param fn DataFunction
**/
func (s *Command) BeforeInsert(fn DataFunctionTx) *Command {
	s.beforeInsert = append(s.beforeInsert, fn)

	return s
}

/**
* BeforeUpdate
* @param fn DataFunction
**/
func (s *Command) BeforeUpdate(fn DataFunctionTx) *Command {
	s.beforeUpdate = append(s.beforeUpdate, fn)

	return s
}

/**
* BeforeDelete
* @param fn Function
**/
func (s *Command) BeforeDelete(fn Function) *Command {
	s.beforeDelete = append(s.beforeDelete, fn)

	return s
}
