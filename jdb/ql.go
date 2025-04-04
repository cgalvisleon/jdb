package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type TypeSelect int

const (
	Select TypeSelect = iota
	Data
)

type Ql struct {
	*QlWhere
	Db         *DB        `json:"-"`
	TypeSelect TypeSelect `json:"type_select"`
	Froms      *QlFroms   `json:"froms"`
	Joins      []*QlJoin  `json:"joins"`
	Selects    []*Field   `json:"selects"`
	Details    []*Field   `json:"details"`
	Rollups    []*Field   `json:"rollups"`
	Groups     []*Field   `json:"group_bys"`
	Havings    *QlHaving  `json:"havings"`
	Orders     *QlOrder   `json:"orders"`
	Sheet      int        `json:"sheet"`
	Offset     int        `json:"offset"`
	Limit      int        `json:"limit"`
	Sql        string     `json:"sql"`
	Result     et.Items   `json:"result"`
}

/**
* Describe
* @return et.Json
**/
func (s *Ql) Describe() et.Json {
	return et.Json{
		"from":     s.listForms(),
		"join":     s.listJoins(),
		"where":    s.listWheres(),
		"group_by": s.listGroups(),
		"having":   s.listHavings(),
		"order_by": s.listOrders(),
		"select":   s.listSelects(),
		"limit":    s.listLimit(),
	}
}

/**
* addFrom
* @param m *Model
* @return *QlFrom
**/
func (s *Ql) addFrom(m *Model) *QlFrom {
	as := string(rune(s.Froms.index))
	from := &QlFrom{
		Model: m,
		As:    as,
	}

	s.Froms.Froms = append(s.Froms.Froms, from)
	s.Froms.index++

	return from
}

/**
* getField
* @param name string bool
* @return *Field
**/
func (s *Ql) getField(name string) *Field {
	findField := func(name string) *Field {
		for _, from := range s.Froms.Froms {
			field := from.getField(name)
			if field != nil {
				field.As = from.As
				return field
			}
		}

		return nil
	}

	for tp, ag := range agregations {
		if ag.re.MatchString(name) {
			n := strs.ReplaceAll(name, []string{ag.Agregation, "(", ")"}, "")
			field := findField(n)
			if field != nil {
				field.Agregation = tp
				return field
			}
		}
	}

	return findField(name)
}

/**
* asField
* @param field *Field
* @return string
**/
func (s *Ql) asField(field *Field) string {
	if len(s.Froms.Froms) <= 1 {
		return field.Name
	}

	return strs.Format("%s.%s", field.Table, field.Name)
}
