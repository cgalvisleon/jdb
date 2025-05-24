package jdb

import (
	"slices"
)

/**
* prepare
* @return *Ql
**/
func (s *Ql) prepare() *Ql {
	if len(s.Selects) == 0 {
		frm := s.Froms.Froms[0]
		for _, col := range frm.Columns {
			if col.Hidden {
				continue
			}

			idx := slices.IndexFunc(s.Hiddens, func(e string) bool { return e == col.Name })
			if idx != -1 {
				continue
			}

			if !slices.Contains([]TypeColumn{TpAtribute}, col.TypeColumn) {
				field := col.GetField()
				field.As = frm.As
				s.setSelect(field)
			}
		}
	} else {
		for _, name := range s.Hiddens {
			idx := slices.IndexFunc(s.Selects, func(e *Field) bool { return e.Name == name })
			if idx != -1 {
				s.Selects = append(s.Selects[:idx], s.Selects[idx+1:]...)
			}

			idx = slices.IndexFunc(s.Details, func(e *Field) bool { return e.Name == name })
			if idx != -1 {
				s.Details = append(s.Details[:idx], s.Details[idx+1:]...)
			}
		}
	}

	s.calcOffset()

	return s
}
