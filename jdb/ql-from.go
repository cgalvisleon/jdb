package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/strs"
)

type QlFrom struct {
	*Model
	As string
}

type QlFroms struct {
	Froms []*QlFrom
	index int
}

func From(m *Model) *Ql {
	result := &Ql{
		Db:         m.Db,
		TypeSelect: Select,
		Froms:      &QlFroms{index: 65, Froms: make([]*QlFrom, 0)},
		Joins:      make([]*QlJoin, 0),
		Selects:    make([]*QlSelect, 0),
		Groups:     make([]*QlSelect, 0),
		Orders:     make([]*QlOrder, 0),
		Details:    make([]*QlSelect, 0),
		Offset:     0,
		Limit:      0,
		Sheet:      0,
	}
	result.QlFilter = &QlFilter{
		main:   result,
		Wheres: make([]*QlWhere, 0),
		Show:   m.Show,
	}
	result.Havings = &QlHaving{
		Ql: result,
	}
	result.Havings.QlFilter = &QlFilter{
		main:   result.Havings,
		Wheres: make([]*QlWhere, 0),
	}

	result.addFrom(m)

	return result
}

/**
* GetField
* @param name string, isCreated bool
* @return *Field
**/
func (s *QlFrom) GetField(name string) *Field {
	result := s.Model.GetField(name)
	if result != nil {
		result.As = s.As
	}

	return result
}

/**
* SetSelectByColumns
* @param selects []*QlSelect, details []*QlSelect
**/
func (s *QlFrom) SetSelectByColumns(selects, details *[]*QlSelect) {
	for _, col := range s.Columns {
		if col.Hidden {
			continue
		}
		if !slices.Contains([]TypeColumn{TpAtribute}, col.TypeColumn) {
			field := col.GetField()
			field.As = s.As
			sel := &QlSelect{
				From:  s,
				Field: field,
			}
			if selects != nil && slices.Contains([]TypeColumn{TpColumn}, col.TypeColumn) {
				*selects = append(*selects, sel)
			} else if details != nil {
				*details = append(*details, sel)
			}
		}
	}

}

/**
* SetSelectBySelects
* @param selects []*QlSelect, details []*QlSelect
**/
func (s *QlFrom) SetSelectBySelects(selects, details *[]*QlSelect) {
	for _, sel := range s.Selects {
		col := sel.Field.Column
		if selects != nil && slices.Contains([]TypeColumn{TpColumn, TpAtribute}, col.TypeColumn) {
			*selects = append(*selects, sel)
		} else if details != nil {
			*details = append(*details, sel)
		}
	}
}

/**
* listForms
* @return []string
**/
func (s *Ql) listForms() []string {
	var result []string
	for _, from := range s.Froms {
		result = append(result, strs.Format(`%s, %s`, from.Table, from.As))
	}

	return result
}
