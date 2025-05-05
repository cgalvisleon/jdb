package jdb

/**
* AfterInsert
* @param fn DataFunction
* @return *Command
**/
func (s *Command) AfterInsert(fn DataFunction) *Command {
	s.afterInsert = append(s.afterInsert, fn)

	return s
}

/**
* AfterUpdate
* @param fn DataFunction
* @return *Command
**/
func (s *Command) AfterUpdate(fn DataFunction) *Command {
	s.afterUpdate = append(s.afterUpdate, fn)

	return s
}

/**
* AfterDelete
* @param fn DataFunction
* @return *Command
**/
func (s *Command) AfterDelete(fn DataFunction) *Command {
	s.afterDelete = append(s.afterDelete, fn)

	return s
}
