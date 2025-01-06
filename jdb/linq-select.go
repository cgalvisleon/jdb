package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/strs"
)

type TypeOperation int

const (
	OperNone TypeOperation = iota
	OpertionSum
	OpertionCount
	OpertionAvg
	OpertionMax
	OpertionMin
)

type LinqSelect struct {
	From      *LinqFrom
	Field     *Field
	Operation TypeOperation
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
				result = append(result, strs.Format(`%s.%s.%s, %s`, frm.As, frm.SourceField.Up(), col.Low(), col.Low()))
			} else {
				result = append(result, strs.Format(`%s.%s, %s`, frm.As, col.Up(), col.Low()))
			}
		}
	}

	return result
}

/**
* GetField
* @param name string
* @return *Field
**/
func (s *LinqSelect) GetField(name string) *Field {
	return s.From.GetField(name)
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
* Return
* @param fields ...string
* @return *Command
**/
func (s *Linq) Return(fields ...string) *Command {
	return nil
}

/**
* Sum
* @param field string
* @return *Linq
**/
func (s *Linq) Sum(field string) *Linq {
	sel := s.GetSelect(field)
	sel.Operation = OpertionSum

	return s
}

/**
* Count
* @param field string
* @return *Linq
**/
func (s *Linq) Count(field string) *Linq {
	sel := s.GetSelect(field)
	sel.Operation = OpertionCount

	return s
}

/**
* Avg
* @param field string
* @return *Linq
**/
func (s *Linq) Avg(field string) *Linq {
	sel := s.GetSelect(field)
	sel.Operation = OpertionAvg

	return s
}

/**
* Max
* @param field string
* @return *Linq
**/
func (s *Linq) Max(field string) *Linq {
	sel := s.GetSelect(field)
	sel.Operation = OpertionMax

	return s
}

/**
* Min
* @param field string
* @return *Linq
**/
func (s *Linq) Min(field string) *Linq {
	sel := s.GetSelect(field)
	sel.Operation = OpertionMin

	return s
}
