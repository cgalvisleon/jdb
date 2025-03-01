package jdb

import (
	"slices"
)

/**
* prepare
* @return *Ql
**/
func (s *Ql) prepare() *Ql {
	isEmpty := len(s.Selects) == 0
	if isEmpty {
		frm := s.Froms.Froms[0]
		for _, col := range frm.Columns {
			if col.Hidden {
				continue
			}

			if !slices.Contains([]TypeColumn{TpAtribute}, col.TypeColumn) {
				field := col.GetField()
				field.As = frm.As
				s.setSelect(field)
			}
		}
	}

	s.calcOffset()

	return s
}
