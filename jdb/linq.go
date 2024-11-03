package jdb

import (
	"strings"

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

type LinqSelect struct {
	From  *LinqFrom
	Field string
}

func (s *LinqSelect) As() string {
	return strs.Format(`%s.%s`, s.From.As, s.Field)
}

type LinqOrder struct {
	LinqSelect
	Sorted bool
}

type Linq struct {
	TypeLinq TypeLinq
	Froms    []*LinqFrom
	Joins    []*LinqJoin
	Wheres   []*LinqWhere
	GroupBys []*LinqSelect
	Havings  []*LinqWhere
	Selects  []*LinqSelect
	Returns  []*LinqSelect
	Orders   []*LinqOrder
	Offset   int
	Limit    int
	index    int
	page     int
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
		if len(list[0]) == 0 {
			return nil
		}

		if len(list[1]) == 0 {
			from := s.Froms[0]
			return &LinqSelect{
				From:  from,
				Field: strs.Uppcase(list[0]),
			}
		}

		if len(list[2]) == 0 {
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
* LinqJoin
* @param m *Model
* @return *Linq
**/
func (s *Linq) Join(m *Model) *LinqJoin {
	return &LinqJoin{
		Linq: s,
	}
}

/**
* Select
* @param columns ...interface{}
* @return *Linq
**/
func (s *Linq) Where(col interface{}) *LinqWhere {
	result := &LinqWhere{
		Linq: s,
	}

	return result
}

/**
* And
* @param col interface{}
* @return *LinqWhere
**/
func (s *Linq) And(col interface{}) *LinqWhere {
	result := &LinqWhere{
		Linq: s,
	}

	return result
}

/**
* And
* @param col interface{}
* @return *LinqWhere
**/
func (s *Linq) Or(col interface{}) *LinqWhere {
	result := &LinqWhere{
		Linq: s,
	}

	return result
}

/**
* GroupBy
* @param columns ...interface{}
* @return *Linq
**/
func (s *Linq) GroupBy(columns ...interface{}) *Linq {
	return s
}

/**
* Having
* @param col interface{}
* @return *LinqWhere
**/
func (s *Linq) Having(col interface{}) *LinqWhere {
	result := &LinqWhere{
		Linq: s,
	}

	return result
}

/**
* Select
* @param columns ...interface{}
* @return *Linq
**/
func (s *Linq) Select(columns ...interface{}) *Linq {
	s.TypeLinq = TypeLinqSelect
	for _, col := range columns {
		c := s.getColumn(col)
		if c != nil {
			s.Selects = append(s.Selects, c)
		}
	}

	return s
}

/**
* Data
* @param columns ...interface{}
* @return *Linq
**/
func (s *Linq) Data(columns ...interface{}) *Linq {
	s.TypeLinq = TypeLinqData
	for _, col := range columns {
		c := s.getColumn(col)
		if c != nil {
			s.Selects = append(s.Selects, c)
		}
	}

	return s
}

/**
* OrderByAsc
* @param columns ...interface{}
* @return *Linq
**/
func (s *Linq) OrderByAsc(columns ...interface{}) *Linq {
	for _, col := range columns {
		c := s.getColumn(col)
		if c != nil {
			order := &LinqOrder{
				LinqSelect: *c,
				Sorted:     true,
			}
			s.Orders = append(s.Orders, order)
		}
	}

	return s
}

/**
* OrderByDesc
* @param columns ...interface{}
* @return *Linq
**/
func (s *Linq) OrderByDesc(columns ...interface{}) *Linq {
	for _, col := range columns {
		c := s.getColumn(col)
		if c != nil {
			order := &LinqOrder{
				LinqSelect: *c,
				Sorted:     false,
			}
			s.Orders = append(s.Orders, order)
		}
	}

	return s
}

/**
* Offset
* @param offset int
* @return *Linq
**/
func (s *Linq) Page(val int) *Linq {
	s.page = val
	s.Offset = s.Limit * (s.page - 1)
	return s
}

/**
* Limit
* @param limit int
* @return *Linq
**/
func (s *Linq) Rows(val int) *Linq {
	s.Limit = val
	s.Offset = s.Limit * (s.page - 1)
	return s
}
