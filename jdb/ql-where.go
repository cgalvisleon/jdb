package jdb

/**
* Where
* @param cond condition
* @return *Cmd
**/
func (s *Ql) Where(cond condition) *Ql {
	s.where.Where(cond)
	return s
}

/**
* And
* @param cond condition
* @return *Cmd
**/
func (s *Ql) And(cond condition) *Ql {
	s.where.And(cond)
	return s
}

/**
* Or
* @param cond condition
* @return *Cmd
**/
func (s *Ql) Or(cond condition) *Ql {
	s.where.Or(cond)
	return s
}
