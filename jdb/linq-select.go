package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type LinqSelect struct {
	From  *LinqFrom
	Field *Field
}

/**
* TypeColumn
* @return TypeColumn
**/
func (s *LinqSelect) TypeColumn() TypeColumn {
	return s.Field.Column.TypeColumn
}

/**
* NewLinqSelect
* @param from *LinqFrom
* @param name string
* @return *LinqSelect
**/
func NewLinqSelect(from *LinqFrom, field *Field) *LinqSelect {
	idx := slices.IndexFunc(from.Selects, func(e *LinqSelect) bool { return e.Field.TableField() == field.TableField() })
	if idx != -1 {
		return from.Selects[idx]
	}

	result := &LinqSelect{
		From:  from,
		Field: field,
	}

	from.Selects = append(from.Selects, result)

	return result
}

/**
* GetField
* @param name string, isCreated bool
* @return *Field
**/
func (s *LinqSelect) GetField(name string, isCreated bool) *Field {
	return s.From.GetField(name, isCreated)
}

/**
* Select
* @param fields ...string
* @return *Linq
**/
func (s *Linq) Select(fields ...string) *Linq {
	for _, field := range fields {
		sel := s.GetSelect(field)
		if sel != nil {
			sel.From.Selects = append(sel.From.Selects, sel)
		}
	}

	return s
}

/**
* Data
* @param fields ...string
* @return *Linq
**/
func (s *Linq) Data(fields ...string) *Linq {
	result := s.Select(fields...)
	result.TypeSelect = Data

	return result
}

/**
* Exec
* @return et.Items, error
**/
func (s *Linq) Exec() (et.Items, error) {
	return et.Items{}, nil
}

/**
* listSelects
* @return []string
**/
func (s *Linq) listSelects() []string {
	result := []string{}
	for _, frm := range s.Froms {
		for _, sel := range frm.Selects {
			result = append(result, strs.Format(`%s, %s`, sel.Field.TableField(), sel.Field.Caption()))
		}
	}

	if len(result) == 0 {
		frm := s.Froms[0]
		for _, col := range frm.Columns {
			if col == frm.SourceField {
				continue
			} else if col.TypeColumn == TpAtribute {
				result = append(result, strs.Format(`%s.%s.%s: %s`, frm.As, frm.SourceField.Up(), col.Low(), col.Low()))
			} else {
				result = append(result, strs.Format(`%s.%s: %s`, frm.As, col.Up(), col.Low()))
			}
		}
	}

	return result
}
