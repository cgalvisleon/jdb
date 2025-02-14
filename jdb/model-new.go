package jdb

import (
	"github.com/cgalvisleon/et/et"
)

/**
* New
* @return et.Json
**/
func (s *Model) New() et.Json {
	var result = et.Json{}
	for _, col := range s.Columns {
		result.Set(col.Name, col.DefaultValue())
	}

	return result
}
