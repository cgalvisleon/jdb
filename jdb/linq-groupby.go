package jdb

import "github.com/cgalvisleon/et/strs"

/**
* GroupBy
* @param columns ...string
* @return *Linq
**/
func (s *Linq) GroupBy(columns ...string) *Linq {
	return s
}

func (s *Linq) ListGroups() []string {
	result := []string{}
	for _, sel := range s.Groups {
		result = append(result, strs.Format(`%s, %s`, sel.Field.Tag(), sel.Field.Caption()))
	}

	return result
}
