package jdb

/**
* Group
* @param fields ...string
* @return *Ql
**/
func (s *Ql) Group(fields ...string) *Ql {
	s.GroupBy = append(s.GroupBy, fields...)
	return s
}
