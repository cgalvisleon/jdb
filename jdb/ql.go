package jdb

import (
	"fmt"
	"strings"
	"sync"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
)

type TypeSelect int

const (
	Select TypeSelect = iota
	Source
)

type Ql struct {
	*QlWhere
	Id         string          `json:"id"`
	Db         *DB             `json:"-"`
	TypeSelect TypeSelect      `json:"type_select"`
	Froms      *QlFroms        `json:"froms"`
	Joins      []*QlJoin       `json:"joins"`
	Selects    []*Field        `json:"selects"`
	Hiddens    []string        `json:"hiddens"`
	Details    []*Field        `json:"details"`
	Groups     []*Field        `json:"group_bys"`
	Havings    *QlHaving       `json:"havings"`
	Orders     *QlOrder        `json:"orders"`
	Concurrent []*Field        `json:"concurrent"`
	Sheet      int             `json:"sheet"`
	Offset     int             `json:"offset"`
	Limit      int             `json:"limit"`
	Sql        string          `json:"sql"`
	Help       et.Json         `json:"help"`
	tx         *Tx             `json:"-"`
	wg         *sync.WaitGroup `json:"-"`
}

/**
* NewQl
* @param db *DB
* @return *Ql
**/
func NewQl(db *DB) *Ql {
	result := &Ql{
		Id:         utility.UUID(),
		Db:         db,
		TypeSelect: Source,
		Froms:      &QlFroms{index: 65, Froms: make([]*QlFrom, 0)},
		Joins:      make([]*QlJoin, 0),
		Selects:    make([]*Field, 0),
		Hiddens:    make([]string, 0),
		Details:    make([]*Field, 0),
		Groups:     make([]*Field, 0),
		Orders:     &QlOrder{Asc: make([]*Field, 0), Desc: make([]*Field, 0)},
		Offset:     0,
		Limit:      0,
		Sheet:      0,
		Help:       helpQl(),
		wg:         &sync.WaitGroup{},
	}
	result.QlWhere = newQlWhere(result.validator)
	result.Havings = NewQlHaving(result)

	return result
}

/**
* Describe
* @return et.Json
**/
func (s *Ql) Describe() et.Json {
	return et.Json{
		"from":     s.getForms(),
		"join":     s.getJoins(),
		"where":    s.getWheres(),
		"group_by": s.getGroupsBy(),
		"having":   s.getHavings(),
		"order_by": s.getOrders(),
		"select":   s.getSelects(),
		"limit":    s.getLimit(),
		"sql":      s.Sql,
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
			field := s.getField(v)
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
		field := s.getField(v)
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

/**
* getColumnField
* @param name string
* @return *Field
**/
func (s *Ql) getColumnField(name string) *Field {
	for _, from := range s.Froms.Froms {
		column := from.getColumn(name)
		if column != nil {
			return column.GetField()
		}
	}

	return newField(name)
}

/**
* getField
* @param name string
* @return *Field
**/
func (s *Ql) getField(name string) *Field {
	findField := func(name string) *Field {
		for _, from := range s.Froms.Froms {
			field := from.getField(name, false)
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

	return fmt.Sprintf("%s.%s", field.Model, field.Name)
}
