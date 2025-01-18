package jdb

import (
	"github.com/cgalvisleon/et/strs"
)

/**
* GroupBy
* @param fields ...string
* @return *Ql
**/
func (s *Ql) GroupBy(fields ...string) *Ql {
	for _, field := range fields {
		sel := s.GetSelect(field)
		if sel != nil {
			s.Groups = append(s.Groups, sel)
		}
	}

	return s
}

/**
* setGroupBy
* @param fields ...string
* @return *Ql
**/
func (s *Ql) setGroupBy(fields ...string) *Ql {
	return s.GroupBy(fields...)
}

/**
* listGroups
* @return []string
**/
func (s *Ql) listGroups() []string {
	result := []string{}
	for _, sel := range s.Groups {
		result = append(result, strs.Format(`%s, %s`, sel.Field.TableField(), sel.Field.Caption()))
	}

	return result
}
