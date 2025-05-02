package jdb

import (
	"github.com/cgalvisleon/et/et"
)

/**
* New
* @return et.Json
**/
func (s *Model) New(fields ...interface{}) et.Json {
	var result = et.Json{}
	for _, col := range s.Columns {
		if s.SourceField != nil && s.SourceField == col {
			continue
		}

		switch col.TypeColumn {
		case TpColumn:
			val := col.DefaultValue()
			result.Set(col.Name, val)
		case TpAtribute:
			val := col.DefaultValue()
			result.Set(col.Name, val)
		case TpCalc:
			for name, fn := range col.CalcFunction {
				val := fn(result)
				result.Set(name, val)
			}
		case TpRelatedTo:
			if col.Detail == nil {
				continue
			}
			with := col.Detail.With
			if with == nil {
				continue
			}
			val := with.New()
			result.Set(col.Name, []et.Json{val})
		case TpRollup:
			if col.Rollup == nil {
				continue
			}
			with := col.Rollup.With
			if with == nil {
				continue
			}
			val := with.New(col.Rollup.Fields...)
			result.Set(col.Name, val)
		}
	}
	return result
}
