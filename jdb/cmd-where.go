package jdb

/**
* Where
* @param cond condition
* @return *Cmd
**/
func (s *Cmd) Where(cond condition) *Cmd {
	s.where.Where(cond)
	return s
}

/**
* And
* @param cond condition
* @return *Cmd
**/
func (s *Cmd) And(cond condition) *Cmd {
	s.where.And(cond)
	return s
}

/**
* Or
* @param cond condition
* @return *Cmd
**/
func (s *Cmd) Or(cond condition) *Cmd {
	s.where.Or(cond)
	return s
}
