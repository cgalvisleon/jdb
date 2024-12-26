package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type Linq struct {
	Db      *DB           `json:"-"`
	Froms   []*LinqFrom   `json:"froms"`
	Joins   []*LinqJoin   `json:"joins"`
	Wheres  []*LinqWhere  `json:"wheres"`
	Groups  []*LinqSelect `json:"group_bys"`
	Havings []*LinqWhere  `json:"havings"`
	Orders  []*LinqOrder  `json:"orders"`
	Sheet   int           `json:"sheet"`
	Offset  int           `json:"offset"`
	Limit   int           `json:"limit"`
	Show    bool          `json:"show"`
	Sql     string        `json:"sql"`
	Result  et.Items      `json:"result"`
	index   int           `json:"-"`
}

func (s *Linq) Describe() et.Json {
	result, err := et.Object(s)
	if err != nil {
		return et.Json{}
	}

	return result
}

func (s *Linq) addFrom(m *Model) *LinqFrom {
	from := &LinqFrom{
		Model:   m,
		As:      string(rune(s.index)),
		Selects: make([]*LinqSelect, 0),
	}

	s.Froms = append(s.Froms, from)
	s.index++

	return from
}

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

func (s *Linq) getSelect(name string) *LinqSelect {
	field := NewField(name)
	if field == nil {
		return nil
	}

	from := s.getFrom(field.TableName())
	if from == nil {
		return nil
	}

	return NewLinqSelect(from, field.Name)
}

/**
* Debug
* @return *Linq
**/
func (s *Linq) Debug() *Linq {
	s.Show = true

	return s
}
