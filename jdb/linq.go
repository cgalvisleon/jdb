package jdb

import (
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type TypeLinq int

const (
	TypeLinqSelect TypeLinq = iota
	TypeLinqData
)

type LinqFrom struct {
	Model
	As string
}

type TypeFuntion int

const (
	Sum TypeFuntion = iota
	Count
	Avg
	Max
	Min
)

type LinqSelect struct {
	From     *LinqFrom
	Field    string
	Function TypeFuntion
}

func (s *LinqSelect) As() string {
	return strs.Format(`%s.%s`, s.From.As, s.Field)
}

type LinqOrder struct {
	LinqSelect
	Sorted bool
}

type Linq struct {
	Db       *DB
	TypeLinq TypeLinq
	Froms    []*LinqFrom
	Joins    []*LinqJoin
	Wheres   []*LinqWhere
	GroupBys []*LinqSelect
	Havings  []*LinqWhere
	Selects  []*LinqSelect
	Orders   []*LinqOrder
	Offset   int
	Limit    int
	Show     bool
	Sql      string
	Result   et.Items
	index    int
	page     int
}

func (s *Linq) Describe() et.Json {
	result, err := et.Object(s)
	if err != nil {
		return et.Json{}
	}

	return result
}

func (s *Linq) addFrom(m Model) *LinqFrom {
	from := &LinqFrom{
		Model: m,
		As:    string(rune(s.index)),
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
		list := strings.Split(v, ".")
		if len(list) == 0 {
			return nil
		}

		if len(list) == 1 {
			for _, from := range s.Froms {
				as := strs.Uppcase(list[0])
				if from.As == as {
					return from
				}
			}

			for _, from := range s.Froms {
				as := strs.Uppcase(list[0])
				if strs.Uppcase(from.Name) == as {
					return from
				}
			}
		}

		if len(list) == 2 {
			table := TableName(list[0], list[1])
			for _, from := range s.Froms {
				if from.Table == table {
					return from
				}
			}
		}

		return nil
	default:
		return nil
	}
}

func (s *Linq) getColumn(col interface{}) *LinqSelect {
	switch v := col.(type) {
	case Column:
		from := s.getFrom(v.Model)
		if from == nil {
			return nil
		}

		return &LinqSelect{
			From:  from,
			Field: v.Field,
		}
	case *Column:
		from := s.getFrom(v.Model)
		if from == nil {
			return nil
		}

		return &LinqSelect{
			From:  from,
			Field: v.Field,
		}
	case string:
		list := strings.Split(v, ".")
		if len(list) == 0 {
			return nil
		}

		if len(list) == 1 {
			from := s.Froms[0]
			return &LinqSelect{
				From:  from,
				Field: strs.Uppcase(list[0]),
			}
		}

		if len(list) == 2 {
			from := s.getFrom(list[0])
			if from == nil {
				return nil
			}

			return &LinqSelect{
				From:  from,
				Field: strs.Uppcase(list[1]),
			}
		}

		if len(list) == 3 {
			from := s.getFrom(list[0] + "." + list[1])
			if from == nil {
				return nil
			}

			return &LinqSelect{
				From:  from,
				Field: strs.Uppcase(list[2]),
			}
		}

		return nil
	default:
		return nil
	}
}

/**
* Debug
* @return *Linq
**/
func (s *Linq) Debug() *Linq {
	s.Show = true

	return s
}
