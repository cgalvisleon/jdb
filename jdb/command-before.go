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
* BeforeInsertOrUpdate
* @param fn DataFunction
**/
func (s *Command) BeforeInsertOrUpdate(fn DataFunctionTx) *Command {
	s.beforeInsert = append(s.beforeInsert, fn)
	s.beforeUpdate = append(s.beforeUpdate, fn)

	return s
}

/**
* BeforeVmInsert
* @param jsCode string
* @return *Command
**/
func (s *Command) BeforeVmInsert(jsCode string) *Command {
	s.beforeVmInsert = append(s.beforeVmInsert, jsCode)

	return s
}

/**
* BeforeVmUpdate
* @param jsCode string
* @return *Command
**/
func (s *Command) BeforeVmUpdate(jsCode string) *Command {
	s.beforeVmUpdate = append(s.beforeVmUpdate, jsCode)

	return s
}

/**
* AfterVmInsert
* @param jsCode string
* @return *Command
**/
func (s *Command) AfterVmInsert(jsCode string) *Command {
	s.afterVmInsert = append(s.afterVmInsert, jsCode)

	return s
}

/**
* AfterVmUpdate
* @param jsCode string
* @return *Command
**/
func (s *Command) AfterVmUpdate(jsCode string) *Command {
	s.afterVmUpdate = append(s.afterVmUpdate, jsCode)

	return s
}

/**
* AfterVmDelete
* @param jsCode string
* @return *Command
**/
func (s *Command) AfterVmDelete(jsCode string) *Command {
	s.afterVmDelete = append(s.afterVmDelete, jsCode)

	return s
}
