package jdb

import (
	"regexp"
)

/**
* Select
* @param fields interface{}
* @return *Ql
**/
func (s *Ql) Select(fields ...string) *Ql {
	if len(s.Froms) == 0 {
		return s
	}

	pattern1 := regexp.MustCompile(`^([A-Za-z0-9]+)\.([A-Za-z0-9]+):([A-Za-z0-9]+)$`)
	pattern2 := regexp.MustCompile(`^([A-Za-z0-9]+)\.([A-Za-z0-9]+)$`)
	pattern3 := regexp.MustCompile(`^([A-Za-z]+)\((.+)\):([A-Za-z0-9]+)$`)
	pattern4 := regexp.MustCompile(`^([A-Za-z]+)\((.+)\)`)
	for _, v := range fields {
		if pattern1.MatchString(v) {
			matches := pattern1.FindStringSubmatch(v)
			if len(matches) == 4 {
				s.Selects[matches[3]] = matches[1]
			}
		} else if pattern2.MatchString(v) {
			matches := pattern2.FindStringSubmatch(v)
			if len(matches) == 3 {
				s.Selects[matches[2]] = matches[1]
			}
		} else if pattern3.MatchString(v) {
			matches := pattern3.FindStringSubmatch(v)
			if len(matches) == 4 {
				s.Selects[matches[3]] = matches[1]
			}
		} else if pattern4.MatchString(v) {
			matches := pattern4.FindStringSubmatch(v)
			if len(matches) == 3 {
				s.Selects[matches[2]] = matches[1]
			}
		} else {

		}

		col, ok := s.Froms.GetColumn(v)
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
		} else if tp == TypeDetail {
			s.Details[v] = s.From.Details[v]
		} else if tp == TypeRollup {
			s.Rollups[v] = s.From.Rollups[v]
		} else if tp == TypeRelation {
			s.Relations[v] = s.From.Relations[v]
		}
	}

	return s
}
