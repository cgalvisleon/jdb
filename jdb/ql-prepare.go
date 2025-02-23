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
			if slices.Contains([]TypeColumn{TpColumn, TpRelatedTo}, col.TypeColumn) {
				field := col.GetField()
				s.setSelect(field)
			}
		}
	}

	return s
}
