package jdb

import (
	"github.com/cgalvisleon/et/et"
)

/**
* New
* @return et.Json
**/
func (s *Model) New(fields ...string) et.Json {
	var result = et.Json{}
	setValue := func(col *Column) {
		switch col.TypeColumn {
		case TpColumn:
			val := col.DefaultValue()
			result.Set(col.Name, val)
		case TpAtribute:
			val := col.DefaultValue()
			result.Set(col.Name, val)
		case TpCalc:
			col.CalcFunction(result)
		case TpRelatedTo:
			if col.Detail == nil {
				return
			}
			with := col.Detail.With
			if with == nil {
				return
			}
			val := with.New()
			result.Set(col.Name, []et.Json{val})
		case TpRollup:
			if col.Rollup == nil {
				return
			}

			rollup := col.Rollup
			with := rollup.With
			if with == nil {
				return
			}
			val := with.New(rollup.Fields...)
			result.Set(col.Name, val)
		}
	}

	if len(fields) == 0 {
		for _, col := range s.Columns {
			if s.SourceField != nil && s.SourceField == col {
				continue
			}

			setValue(col)
		}
	} else {
		for _, field := range fields {
			col := s.getColumn(field)
			if col == nil {
				continue
			}
			setValue(col)
		}
	}

	return result
}
