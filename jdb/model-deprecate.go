package jdb

import "github.com/cgalvisleon/et/et"

/**
* setFields
* @param fields et.Json
**/
func (s *Model) setFields(fields et.Json) {
	s.defineFields(fields)
}
