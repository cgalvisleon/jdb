package jdb

/**
* Where
* @param cond Condition
* @return *Cmd
**/
func (s *Ql) Where(cond Condition) *Ql {
	s.where.Where(cond)
	return s
}

/**
* And
* @param cond Condition
* @return *Cmd
**/
func (s *Ql) And(cond Condition) *Ql {
	s.where.And(cond)
	return s
}

/**
* Or
* @param cond Condition
* @return *Cmd
**/
func (s *Ql) Or(cond Condition) *Ql {
	s.where.Or(cond)
	return s
}
