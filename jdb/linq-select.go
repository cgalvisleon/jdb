package jdb

import "github.com/cgalvisleon/et/strs"

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
* NewLinqSelect
* @param from *LinqFrom
* @param name string
* @return *LinqSelect
**/
func NewLinqSelect(from *LinqFrom, name string) *LinqSelect {
	field := from.GetField(name)
	if field == nil {
		return nil
	}

	result := &LinqSelect{
		From:  from,
		Field: field,
	}

	from.Selects = append(from.Selects, result)

	return result
}

/**
* GetSelects
* @return []string
**/
func (s *Linq) ListSelects() []string {
	result := []string{}
	for _, frm := range s.Froms {
		for _, sel := range frm.Selects {
			result = append(result, strs.Format(`%s, %s`, sel.Field.Tag(), sel.Field.Caption()))
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
	for _, name := range fields {
		sel := s.getSelect(name)
		if sel != nil {
			sel.From.Selects = append(sel.From.Selects, sel)
		}
	}

	return s
}

/**
* Sum
* @param field string
* @return *Linq
**/
func (s *Linq) Sum(field string) *Linq {
	sel := s.getSelect(field)
	sel.Operation = OpertionSum

	return s
}

/**
* Count
* @param field string
* @return *Linq
**/
func (s *Linq) Count(field string) *Linq {
	sel := s.getSelect(field)
	sel.Operation = OpertionCount

	return s
}

/**
* Avg
* @param field string
* @return *Linq
**/
func (s *Linq) Avg(field string) *Linq {
	sel := s.getSelect(field)
	sel.Operation = OpertionAvg

	return s
}

/**
* Max
* @param column string
* @return *Linq
**/
func (s *Linq) Max(column string) *Linq {
	sel := s.getSelect(column)
	sel.Operation = OpertionMax

	return s
}

/**
* Min
* @param column string
* @return *Linq
**/
func (s *Linq) Min(column string) *Linq {
	sel := s.getSelect(column)
	sel.Operation = OpertionMin

	return s
}
