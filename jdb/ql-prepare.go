package jdb

/**
* prepare
* @return *Ql
**/
func (s *Ql) prepare() *Ql {
	isEmpty := len(s.Selects) == 0
	if isEmpty {
		frm := s.Froms.Froms[0]
		for _, col := range frm.Columns {
			s.setSelect(col.GetField())
		}
	}

	return s
}
