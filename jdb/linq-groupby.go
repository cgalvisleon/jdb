package jdb

import (
	"github.com/cgalvisleon/et/strs"
)

/**
* GroupBy
* @param fields ...string
* @return *Linq
**/
func (s *Linq) GroupBy(fields ...string) *Linq {
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
* @return *Linq
**/
func (s *Linq) setGroupBy(fields ...string) *Linq {
	return s.GroupBy(fields...)
}

/**
* listGroups
* @return []string
**/
func (s *Linq) listGroups() []string {
	result := []string{}
	for _, sel := range s.Groups {
		result = append(result, strs.Format(`%s, %s`, sel.Field.TableField(), sel.Field.Caption()))
	}

	return result
}
