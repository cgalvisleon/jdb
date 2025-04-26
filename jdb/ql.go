package jdb

import (
	"strings"

	"github.com/cgalvisleon/et/console"
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
	Help       et.Json    `json:"help"`
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
		"help":     s.Help,
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
* validator
* validate this val is a field or basic type
* @param val interface{}
* @return interface{}
**/
func (s *Ql) validator(val interface{}) interface{} {
	console.Debug("validator:", val)
	switch v := val.(type) {
	case string:
		if strings.HasPrefix(v, "$") {
			v = strings.TrimPrefix(v, "$")
			field := s.getField(v, false)
			if field != nil {
				return field
			}
			return nil
		}

		v = strings.Replace(v, `\\$`, `\$`, 1)
		v = strings.Replace(v, `\$`, `$`, 1)
		field := s.getField(v, false)
		if field != nil {
			return field
		}

		return v
	case *Field:
		return v
	case Field:
		return v
	case *Column:
		return v.GetField()
	case Column:
		return v.GetField()
	default:
		return v
	}
}

func (s *Ql) getColumnField(name string, isCreated bool) *Field {
	for _, from := range s.Froms.Froms {
		column := from.getColumn(name)
		if column != nil {
			return column.GetField()
		}
	}

	if isCreated {
		return &Field{
			Name: name,
		}
	}

	return nil
}

/**
* getField
* @param name string, isCreated bool
* @return *Field
**/
func (s *Ql) getField(name string, isCreated bool) *Field {
	findField := func(name string) *Field {
		for _, from := range s.Froms.Froms {
			field := from.getField(name, isCreated)
			if field != nil {
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
