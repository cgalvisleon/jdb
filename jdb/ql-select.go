package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type QlSelect struct {
	From  *QlFrom
	Field *Field
}

/**
* TypeColumn
* @return TypeColumn
**/
func (s *QlSelect) TypeColumn() TypeColumn {
	return s.Field.Column.TypeColumn
}

/**
* NewQlSelect
* @param from *QlFrom
* @param name string
* @return *QlSelect
**/
func NewQlSelect(from *QlFrom, field *Field) *QlSelect {
	idx := slices.IndexFunc(from.Selects, func(e *QlSelect) bool { return e.Field.TableField() == field.TableField() })
	if idx != -1 {
		return from.Selects[idx]
	}

	result := &QlSelect{
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
func (s *QlSelect) GetField(name string, isCreated bool) *Field {
	return s.From.GetField(name, isCreated)
}

/**
* Select
* @param fields ...string
* @return *Ql
**/
func (s *Ql) Select(fields ...string) *Ql {
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
* @return *Ql
**/
func (s *Ql) Data(fields ...string) *Ql {
	result := s.Select(fields...)
	result.TypeSelect = Data

	return result
}

/**
* Exec
* @return et.Items, error
**/
func (s *Ql) Exec() (et.Items, error) {
	return et.Items{}, nil
}

/**
* listSelects
* @return []string
**/
func (s *Ql) listSelects() []string {
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
				result = append(result, strs.Format(`%s.%s.%s: %s`, frm.As, frm.SourceField.Name, col.Name, col.Name))
			} else {
				result = append(result, strs.Format(`%s.%s: %s`, frm.As, col.Name, col.Name))
			}
		}
	}

	return result
}
