package jdb

/**
* Where
* @param cond Condition
* @return *Cmd
**/
func (s *Cmd) Where(cond Condition) *Cmd {
	s.where.Where(cond)
	return s
}

/**
* And
* @param cond Condition
* @return *Cmd
**/
func (s *Cmd) And(cond Condition) *Cmd {
	s.where.And(cond)
	return s
}

/**
* Or
* @param cond Condition
* @return *Cmd
**/
func (s *Cmd) Or(cond Condition) *Cmd {
	s.where.Or(cond)
	return s
}
