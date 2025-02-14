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
	Groups     []*Field   `json:"group_bys"`
	Havings    *QlHaving  `json:"havings"`
	Orders     []*QlOrder `json:"orders"`
	Sheet      int        `json:"sheet"`
	Offset     int        `json:"offset"`
	Limit      int        `json:"limit"`
	Sql        string     `json:"sql"`
	Result     et.Items   `json:"result"`
	index      int        `json:"-"`
}

/**
* Describe
* @return et.Json
**/
func (s *Ql) Describe() et.Json {
	result, err := et.Object(s)
	if err != nil {
		return et.Json{}
	}

	return result
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
	s.index++

	return from
}

/**
* getFrom
* @param v string
* @return *QlFrom
**/
func (s *Ql) getFrom(v string) *QlFrom {
	for _, from := range s.Froms.Froms {
		if from.As == v && from.Model.Name == v {
			return from
		}
	}

	return nil
}

/**
* GetField
* @param name string bool
* @return *Field
**/
func (s *Ql) GetField(name string) *Field {
	var field *Field
	for _, from := range s.Froms.Froms {
		field = from.GetField(name)
		if field != nil {
			field.As = from.As
			return field
		}
	}

	return nil
}

/**
* GetAgregation
* @params name string
* @return *Field
**/
func (s *Ql) GetAgregation(name string) *Field {
	for tp, ag := range agregations {
		if ag.re.MatchString(name) {
			name = strs.ReplaceAll(name, []string{ag.Agregation, "(", ")"}, "")
			field := s.GetField(name)
			if field != nil {
				field.Agregation = tp
				return field
			}
		}
	}

	return nil
}
