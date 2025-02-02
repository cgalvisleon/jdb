package jdb

/**
* prepare
* @return *Ql
**/
func (s *Ql) prepare() *Ql {
	for _, frm := range s.Froms {
		frm.SetSelectBySelects(&s.Selects, &s.Details)
	}

	isEmpty := len(s.Selects)+len(s.Details) == 0
	if isEmpty {
		frm := s.Froms[0]
		frm.SetSelectByColumns(&s.Selects, &s.Details)
	}

	return s
}
