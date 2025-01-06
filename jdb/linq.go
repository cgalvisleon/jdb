package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
)

type TypeSelect int

const (
	Select TypeSelect = iota
	Data
)

type Linq struct {
	Db         *DB           `json:"-"`
	TypeSelect TypeSelect    `json:"type_select"`
	Froms      []*LinqFrom   `json:"froms"`
	Joins      []*LinqJoin   `json:"joins"`
	Wheres     []*LinqWhere  `json:"wheres"`
	Groups     []*LinqSelect `json:"group_bys"`
	Havings    []*LinqWhere  `json:"havings"`
	Orders     []*LinqOrder  `json:"orders"`
	Sheet      int           `json:"sheet"`
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
	Show       bool          `json:"show"`
	Sql        string        `json:"sql"`
	Result     et.Items      `json:"result"`
	index      int           `json:"-"`
}

func (s *Linq) Describe() et.Json {
	result, err := et.Object(s)
	if err != nil {
		return et.Json{}
	}

	return result
}

/**
* addFrom
* @param m *Model
* @return *LinqFrom
**/
func (s *Linq) addFrom(m *Model) *LinqFrom {
	as := string(rune(s.index))
	from := &LinqFrom{
		Model:   m,
		As:      as,
		Selects: make([]*LinqSelect, 0),
	}

	s.Froms = append(s.Froms, from)
	s.index++

	return from
}

/**
* getFrom
* @param m interface{}
* @return *LinqFrom
**/
func (s *Linq) getFrom(m interface{}) *LinqFrom {
	switch v := m.(type) {
	case Model:
		for _, from := range s.Froms {
			if from.Table == v.Table {
				return from
			}
		}

		return nil
	case *Model:
		for _, from := range s.Froms {
			if from.Table == v.Table {
				return from
			}
		}

		return nil
	case string:
		for _, from := range s.Froms {
			if from.Table == strs.Lowcase(v) {
				return from
			}
		}

		return nil
	default:
		return nil
	}
}

/**
* GetField
* @param name string
* @return *LinqSelect
**/
func (s *Linq) GetField(name string) *Field {
	var field *Field
	for _, from := range s.Froms {
		field = from.GetField(name)
		if field != nil {
			field.Owner = from
			break
		}
	}

	return field
}

/**
* GetSelect
* @param name string
* @return *LinqSelect
**/
func (s *Linq) GetSelect(name string) *LinqSelect {
	field := s.GetField(name)

	if field != nil {
		return NewLinqSelect(field.Owner.(*LinqFrom), field)
	}

	return nil
}

/**
* Debug
* @return *Linq
**/
func (s *Linq) Debug() *Linq {
	s.Show = true

	return s
}

/**
* Query
* @param query et.Json
* @return et.Items, error
**/
func Query(query et.Json) (et.Items, error) {
	if query.IsEmpty() {
		return et.Items{}, mistake.New(MSG_QUERY_EMPTY)
	}

	from := query.Str("from")
	if !utility.ValidStr(from, 0, []string{""}) {
		return et.Items{}, mistake.New(MSG_QUERY_FROM_REQUIRED)
	}

	model := models[from]
	if model == nil {
		return et.Items{}, mistake.Newf(MSG_MODEL_NOT_FOUND, from)
	}

	return From(model).
		Debug().
		Query(query)
}
