package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/et"
)

/**
* New
* @return et.Json
**/
func (s *Model) New() et.Json {
	var result = et.Json{}
	for _, col := range s.Columns {
		if !slices.Contains([]TypeColumn{TpColumn, TpAtribute}, col.TypeColumn) {
			continue
		}

		if s.SourceField != nil && s.SourceField == col {
			continue
		}

		result.Set(col.Name, col.DefaultValue())
	}

	return result
}
