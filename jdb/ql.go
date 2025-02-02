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

type Ql struct {
	*QlFilter
	Db         *DB              `json:"-"`
	TypeSelect TypeSelect       `json:"type_select"`
	Froms      []*QlFrom        `json:"froms"`
	Selects    []*QlSelect      `json:"selects"`
	Joins      []*QlJoin        `json:"joins"`
	Groups     []*QlSelect      `json:"group_bys"`
	Havings    *QlHaving        `json:"havings"`
	Orders     []*QlOrder       `json:"orders"`
	Details    []*QlSelect      `json:"details"`
	Generateds []*FuncGenerated `json:"generateds"`
	Sheet      int              `json:"sheet"`
	Offset     int              `json:"offset"`
	Limit      int              `json:"limit"`
	Sql        string           `json:"sql"`
	Source     string           `json:"source"`
	Result     et.Items         `json:"result"`
	index      int              `json:"-"`
}

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
	as := string(rune(s.index))
	from := &QlFrom{
		Model:   m,
		As:      as,
		Selects: make([]*QlSelect, 0),
	}

	s.Source = m.Source
	s.Froms = append(s.Froms, from)
	s.index++

	return from
}

/**
* getFrom
* @param m interface{}
* @return *QlFrom
**/
func (s *Ql) getFrom(m interface{}) *QlFrom {
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
* GetAgregation
* @params name string
* @return *Field
**/
func (s *Ql) GetAgregation(name string) *Field {
	for tp, ag := range agregations {
		if ag.re.MatchString(name) {
			name = strs.ReplaceAll(name, []string{ag.Agregation, "(", ")"}, "")
			field := s.GetField(name, false)
			if field != nil {
				field.Agregation = tp
				return field
			}
		}
	}

	return nil
}

/**
* GetField
* @param name string, isCreated bool
* @return *Field
**/
func (s *Ql) GetField(name string, isCreated bool) *Field {
	var field *Field
	for _, from := range s.Froms {
		field = from.GetField(name, isCreated)
		if field != nil {
			field.Owner = from
			return field
		}
	}

	return nil
}

/**
* GetSelect
* @param name string
* @return *QlSelect
*
 */
func (s *Ql) GetSelect(name string) *QlSelect {
	field := s.GetAgregation(name)
	if field != nil {
		return NewQlSelect(field.Owner.(*QlFrom), field)
	}

	field = s.GetField(name, true)
	if field == nil {
		return nil
	}

	if field.Column.TypeColumn == TpDetail {
		col := field.Column
		model := col.Model
		detail := col.Detail
		fk := strs.Format(`%s.%s`, model.Name, col.Model.KeyField.Name)
		s.Join(detail).
			On(model.Name).Eq(fk)
	}

	return NewQlSelect(field.Owner.(*QlFrom), field)
}

/**
* Debug
* @return *Ql
**/
func (s *Ql) Debug() *Ql {
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

/**
* Commands
* @param command et.Json
* @return et.Items, error
**/
func Commands(command et.Json) (et.Items, error) {
	if command.IsEmpty() {
		return et.Items{}, mistake.New(MSG_QUERY_EMPTY)
	}

	from := command.Str("from")
	if !utility.ValidStr(from, 0, []string{""}) {
		return et.Items{}, mistake.New(MSG_QUERY_FROM_REQUIRED)
	}

	model := models[from]
	if model == nil {
		return et.Items{}, mistake.Newf(MSG_MODEL_NOT_FOUND, from)
	}

	return model.
		Command(command)
}
