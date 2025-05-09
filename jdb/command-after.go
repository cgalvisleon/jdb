package jdb

/**
* AfterInsert
* @param fn DataFunction
* @return *Command
**/
func (s *Command) AfterInsert(fn DataFunctionTx) *Command {
	s.afterInsert = append(s.afterInsert, fn)

	return s
}

/**
* AfterUpdate
* @param fn DataFunctionTx
* @return *Command
**/
func (s *Command) AfterUpdate(fn DataFunctionTx) *Command {
	s.afterUpdate = append(s.afterUpdate, fn)

	return s
}

/**
* AfterDelete
* @param fn Function
* @return *Command
**/
func (s *Command) AfterDelete(fn Function) *Command {
	s.afterDelete = append(s.afterDelete, fn)

	return s
}
