package jdb

/**
* AfterInsert
* @param fn DataFunction
* @return *Command
**/
func (s *Model) AfterInsert(fn DataFunctionTx) *Model {
	s.afterInsert = append(s.afterInsert, fn)

	return s
}

/**
* AfterUpdate
* @param fn DataFunctionTx
* @return *Command
**/
func (s *Model) AfterUpdate(fn DataFunctionTx) *Model {
	s.afterUpdate = append(s.afterUpdate, fn)

	return s
}

/**
* AfterDelete
* @param fn DataFunctionTx
* @return *Command
**/
func (s *Model) AfterDelete(fn DataFunctionTx) *Model {
	s.afterDelete = append(s.afterDelete, fn)

	return s
}

/**
* AfterInsertOrUpdate
* @param fn DataFunctionTx
* @return *Command
**/
func (s *Model) AfterInsertOrUpdate(fn DataFunctionTx) *Model {
	s.afterInsert = append(s.afterInsert, fn)
	s.afterUpdate = append(s.afterUpdate, fn)

	return s
}
