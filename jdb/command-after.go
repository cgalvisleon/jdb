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
* @param fn DataFunctionTx
* @return *Command
**/
func (s *Command) AfterDelete(fn DataFunctionTx) *Command {
	s.afterDelete = append(s.afterDelete, fn)

	return s
}

/**
* AfterInsertOrUpdate
* @param fn DataFunctionTx
* @return *Command
**/
func (s *Command) AfterInsertOrUpdate(fn DataFunctionTx) *Command {
	s.afterInsert = append(s.afterInsert, fn)
	s.afterUpdate = append(s.afterUpdate, fn)

	return s
}

/**
* AfterFuncInsert
* @param jsCode string
* @return *Command
**/
func (s *Command) AfterFuncInsert(jsCode string) *Command {
	s.afterFuncInsert = append(s.afterFuncInsert, jsCode)

	return s
}

/**
* AfterFuncUpdate
* @param jsCode string
* @return *Command
**/
func (s *Command) AfterFuncUpdate(jsCode string) *Command {
	s.afterFuncUpdate = append(s.afterFuncUpdate, jsCode)

	return s
}

/**
* AfterVmDelete
* @param jsCode string
* @return *Command
**/
func (s *Command) AfterFuncDelete(jsCode string) *Command {
	s.afterFuncDelete = append(s.afterFuncDelete, jsCode)

	return s
}

/**
* AfterInsertOrUpdateFunc
* @param tp TypeEvent, jsCode string
* @return *Command
**/
func (s *Command) AfterInsertOrUpdateFunc(tp TypeEvent, jsCode string) *Command {
	switch tp {
	case EventInsert:
		s.afterFuncInsert = append(s.afterFuncInsert, jsCode)
	case EventUpdate:
		s.afterFuncUpdate = append(s.afterFuncUpdate, jsCode)
	}
	return s
}
