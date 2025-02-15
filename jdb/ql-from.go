package jdb

import (
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
		QlWhere:    NewQlWhere(),
		Selects:    make([]*Field, 0),
		Details:    make([]*Field, 0),
		Groups:     make([]*Field, 0),
		Orders:     &QlOrder{Asc: make([]*Field, 0), Desc: make([]*Field, 0)},
		Offset:     0,
		Limit:      0,
		Sheet:      0,
	}
	result.Havings = &QlHaving{Ql: result, QlWhere: NewQlWhere()}
	result.addFrom(m)

	return result
}

/**
* getField
* @param name string, isCreated bool
* @return *Field
**/
func (s *QlFrom) getField(name string) *Field {
	result := s.Model.GetField(name)
	if result != nil {
		result.As = s.As
	}

	return result
}

/**
* listForms
* @return []string
**/
func (s *Ql) listForms() []string {
	var result []string
	for _, from := range s.Froms.Froms {
		result = append(result, strs.Format(`%s, %s`, from.Table, from.As))
	}

	return result
}
