package jdb

import (
	"fmt"
	"strings"
)

/**
* Select
* @param fields interface{}
* @return *Ql
**/
func (s *Ql) Select(fields ...string) *Ql {
	if s.From == nil {
		return s
	}

	isDotInMiddle := func(s string) bool {
		if s == "" {
			return false
		}
		dot := strings.Index(s, ".")
		if dot <= 0 || dot == len(s)-1 {
			return false
		}
		return strings.LastIndex(s, ".") == dot // asegura que hay solo un punto
	}

	for _, v := range fields {
		if isDotInMiddle(v) {
			as := strings.Split(v, ".")[1]
			s.Selects[v] = as
			continue
		}

		col, ok := s.From.GetColumn(v)
		tp := col.String("type")
		if ok && TypeColumn[tp] {
			s.Selects[v] = v
			continue
		}

		if s.From.UseAtribs() || TypeAtrib[tp] {
			s.Atribs[v] = v
			continue
		}

		if tp == TypeCalc {
			s.Calcs[v] = s.From.Calcs[v]
		} else if tp == TypeVm {
			s.Vms[v] = s.From.Vms[v]
		} else if tp == TypeRollup {
			s.Rollups[v] = s.From.Rollups[v]
		} else if tp == TypeRelation {
			s.Relations[v] = s.From.Relations[v]
		} else if tp == TypeDetail {
			to, err := s.db.GetModel(v)
			if err != nil {
				continue
			}

			detail := s.From.Details[v]
			references := detail.Json("references")
			columns := references.ArrayJson("columns")
			as := string(rune(len(s.Joins) + 66))
			first := true
			for _, fk := range columns {
				for k, v := range fk {
					if first {
						s.Join(to, as, Eq(fmt.Sprintf("A.%s", k), fmt.Sprintf("%s.%s", as, v)))
						first = false
						continue
					}
					s.And(Eq(fmt.Sprintf("A.%s", k), fmt.Sprintf("%s.%s", as, v)))
				}
			}
		}
	}

	return s
}
