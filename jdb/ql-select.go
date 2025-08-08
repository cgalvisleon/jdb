package jdb

import (
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
		idx := slices.IndexFunc(s.Selects, func(e *Field) bool { return e == field })
		if idx == -1 {
			s.Selects = append(s.Selects, field)
		}

		if field.Column.CalcFunction != nil {
			s.Details = append(s.Details, field)
		}
	} else {
		if slices.Contains([]TypeColumn{TpRollup}, field.Column.TypeColumn) {
			rollup := field.Column.Rollup
			for _, name := range rollup.Fields {
				field.Select = append(field.Select, name)
			}
		}

		idx := slices.IndexFunc(s.Details, func(e *Field) bool { return e == field })
		if idx == -1 {
			s.Details = append(s.Details, field)
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
			field := s.getField(key)
			if field.Column.TypeColumn == TpRelatedTo {
				s.setDetail(v)
			}
		}
	}

	for _, name := range fields {
		switch v := name.(type) {
		case string:
			field := s.getField(v)
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
	result.TypeSelect = Source

	return result
}

/**
* Hidden
* @param fields ...string
* @return *Ql
**/
func (s *Ql) Hidden(fields ...string) *Ql {
	return s.setHidden(fields...)
}

/**
* setSelects
* @param fields ...interface{}
* @return *Ql
**/
func (s *Ql) SetSelects(fields ...interface{}) *Ql {
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
* setHidden
* @param columns ...*Column
* @return *Ql
**/
func (s *Ql) setHidden(columns ...string) *Ql {
	s.Hiddens = append(s.Hiddens, columns...)

	return s
}

/**
* getSelects
* @return []string
**/
func (s *Ql) getSelects() []string {
	result := []string{}
	for _, sel := range s.Selects {
		result = append(result, sel.asField())
	}

	return result
}
