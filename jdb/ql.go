package jdb

import (
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type TypeSelect int

const (
	Select TypeSelect = iota
	Source
)

type Ql struct {
	*QlWhere
	Id         string     `json:"id"`
	Db         *DB        `json:"-"`
	TypeSelect TypeSelect `json:"type_select"`
	Froms      *QlFroms   `json:"froms"`
	Joins      []*QlJoin  `json:"joins"`
	Selects    []*Field   `json:"selects"`
	Hiddens    []string   `json:"hiddens"`
	Details    []*Field   `json:"details"`
	Groups     []*Field   `json:"group_bys"`
	Havings    *QlHaving  `json:"havings"`
	Orders     *QlOrder   `json:"orders"`
	Sheet      int        `json:"sheet"`
	Offset     int        `json:"offset"`
	Limit      int        `json:"limit"`
	Sql        string     `json:"sql"`
	Help       et.Json    `json:"help"`
	tx         *Tx        `json:"-"`
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
* setTx
* @param tx *Tx
* @return *Ql
**/
func (s *Ql) setTx(tx *Tx) *Ql {
	s.tx = tx

	return s
}

/**
* Tx
* @return *Tx
**/
func (s *Ql) Tx() *Tx {
	return s.tx
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
	switch v := val.(type) {
	case string:
		if strings.HasPrefix(v, ":") {
			v = strings.TrimPrefix(v, ":")
			field := s.getField(v, false)
			if field != nil {
				return field
			}
			return nil
		}

		if strings.HasPrefix(v, "$") {
			v = strings.TrimPrefix(v, "$")
			return v
		}

		v = strings.Replace(v, `\\:`, `\:`, 1)
		v = strings.Replace(v, `\:`, `:`, 1)
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
	case []interface{}:
		return v
	case []string:
		return v
	case []et.Json:
		return v
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
		return newField(name)
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

	return strs.Format("%s.%s", field.Model, field.Name)
}
