package jdb

import (
	"fmt"
	"slices"

	"github.com/cgalvisleon/et/et"
)

/**
* setSelect
* @param field *Field
* @return *Ql
**/
func (s *Ql) setSelect(field *Field) *Ql {
	if field == nil || field.Column == nil {
		return s
	}

	if slices.Contains([]TypeColumn{TpColumn, TpAtribute}, field.Column.TypeColumn) {
		idx := slices.IndexFunc(s.Selects, func(e *Field) bool { return e.asField() == field.asField() })
		if idx == -1 {
			s.Selects = append(s.Selects, field)
		}
	} else if slices.Contains([]TypeColumn{TpCalc, TpRelatedTo}, field.Column.TypeColumn) {
		idx := slices.IndexFunc(s.Details, func(e *Field) bool { return e.asField() == field.asField() })
		if idx == -1 {
			s.Details = append(s.Details, field)
		}
	} else if slices.Contains([]TypeColumn{TpRollup}, field.Column.TypeColumn) {
		name := field.Column.Rollup.Fields[field.Name]
		def := fmt.Sprintf(`%s:%s`, name, field.Name)
		idx := slices.IndexFunc(s.Details, func(e *Field) bool { return e.Column.Rollup == field.Column.Rollup })
		if idx == -1 {
			field.Select = append(field.Select, def)
			s.Details = append(s.Details, field)
		} else {
			s.Details[idx].Select = append(s.Details[idx].Select, def)
		}
	}

	return s
}

/**
* Select
* @param fields ...interface{}
* @return *Ql
**/
func (s *Ql) Select(fields ...interface{}) *Ql {
	setRelationTo := func(v map[string]interface{}) {
		for key := range v {
			field := s.getField(key, true)
			if field.Column.TypeColumn == TpRelatedTo {
				s.setDetail(v)
			}
		}
	}

	for _, name := range fields {
		switch v := name.(type) {
		case string:
			field := s.getField(v, true)
			s.setSelect(field)
		case et.Json:
			setRelationTo(v)
		case map[string]interface{}:
			setRelationTo(v)
		}
	}
	s.TypeSelect = Select

	return s
}

/**
* Data
* @param fields ...interface{}
* @return *Ql
**/
func (s *Ql) Data(fields ...interface{}) *Ql {
	result := s.Select(fields...)
	result.TypeSelect = Data

	return result
}

/**
* setSelects
* @param fields ...interface{}
* @return *Ql
**/
func (s *Ql) setSelects(fields ...interface{}) *Ql {
	froms := s.Froms.Froms
	if len(froms) == 0 {
		return s
	}

	model := froms[0].Model
	if model == nil {
		return s
	}

	if model.SourceField != nil {
		s.Data(fields...)
	} else {
		s.Select(fields...)
	}

	return s
}

/**
* ListSelects
* @return []string
**/
func (s *Ql) listSelects() []string {
	result := []string{}
	for _, sel := range s.Selects {
		result = append(result, sel.asField())
	}

	return result
}
