package jdb

import "github.com/cgalvisleon/et/et"

type FilterTo interface {
	And(val interface{}) *LinqFilter
	Or(val interface{}) *LinqFilter
	Select(fields ...string) *Linq
	Data(fields ...string) *Linq
	Exec() (et.Items, error)
	One() (et.Item, error)
}

type LinqFilter struct {
	main   FilterTo
	where  *LinqWhere
	Wheres []*LinqWhere
}

func (s *LinqFilter) setCondition(where et.Json) *LinqFilter {
	if where["eq"] != nil {
		s.Eq(where["eq"])
	} else if where["neg"] != nil {
		s.Neg(where["neg"])
	} else if where["in"] != nil {
		s.In(where["in"])
	} else if where["like"] != nil {
		s.Like(where["like"])
	} else if where["more"] != nil {
		s.More(where["more"])
	} else if where["less"] != nil {
		s.Less(where["less"])
	} else if where["moreEq"] != nil {
		s.MoreEq(where["moreEq"])
	} else if where["lessEq"] != nil {
		s.LessEs(where["lessEq"])
	} else if where["search"] != nil {
		s.Search(where["search"])
	} else if where["between"] != nil {
		s.Between(where["between"])
	} else if where["isNull"] != nil {
		s.IsNull()
	} else if where["notNull"] != nil {
		s.NotNull()
	}

	return s
}

/**
* AddValue
* @param val interface{}
* @return *LinqFilter
**/
func (s *LinqFilter) AddValue(val interface{}) FilterTo {
	appendValue := func(linq *Linq, value interface{}) {
		switch v := value.(type) {
		case string:
			field := linq.GetField(v, false)
			if field != nil {
				s.where.Values = append(s.where.Values, field)
			} else {
				s.where.Values = append(s.where.Values, value)
			}
		default:
			s.where.Values = append(s.where.Values, value)
		}
	}

	switch m := s.main.(type) {
	case *Linq:
		appendValue(m, val)
	case *Command:
		s.where.Values = append(s.where.Values, val)
	case *LinqJoin:
		appendValue(m.Linq, val)
	case *LinqHaving:
		appendValue(m.Linq, val)
	}

	s.Wheres = append(s.Wheres, s.where)
	return s.main
}

/**
* Eq
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) Eq(val interface{}) FilterTo {
	s.where.Operator = Equal
	return s.AddValue(val)
}

/**
* Neg
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) Neg(val interface{}) FilterTo {
	s.where.Operator = Neg
	return s.AddValue(val)
}

/**
* In
* @param val ...any
* @return FilterTo
**/
func (s *LinqFilter) In(val ...any) FilterTo {
	s.where.Operator = In
	return s.AddValue(val)
}

/**
* Like
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) Like(val interface{}) FilterTo {
	s.where.Operator = Like
	return s.AddValue(val)
}

/**
* More
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) More(val interface{}) FilterTo {
	s.where.Operator = More
	return s.AddValue(val)
}

/**
* Less
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) Less(val interface{}) FilterTo {
	s.where.Operator = Less
	return s.AddValue(val)
}

/**
* MoreEq
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) MoreEq(val interface{}) FilterTo {
	s.where.Operator = MoreEq
	return s.AddValue(val)
}

/**
* LessEq
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) LessEs(val interface{}) FilterTo {
	s.where.Operator = LessEq
	return s.AddValue(val)
}

/**
* Search
* @param val interface{}
* @return FilterTo
**/
func (s *LinqFilter) Search(val interface{}) FilterTo {
	s.where.Operator = Search
	return s.AddValue(val)
}

/**
* Between
* @param val1, val2 interface{}
* @return FilterTo
**/
func (s *LinqFilter) Between(val interface{}) FilterTo {
	s.where.Operator = Between
	vals, ok := val.([]interface{})
	if !ok {
		return s.main
	}

	switch len(vals) {
	case 1:
		return s.AddValue(vals[0])
	case 2:
		s.AddValue(vals[0])
		return s.AddValue(vals[1])
	default:
		return s.main
	}
}

/**
* IsNull
* @return *LinqFilter
**/
func (s *LinqFilter) IsNull() *LinqFilter {
	s.where.Operator = IsNull
	return s
}

/**
* NotNull
* @return *LinqFilter
**/
func (s *LinqFilter) NotNull() *LinqFilter {
	s.where.Operator = NotNull
	return s
}

/**
* First
* @param n int
* @return et.Items, error
**/
func (s *LinqFilter) First(n int) (et.Items, error) {
	return s.main.(*Linq).First(n)
}

/**
* All
* @return et.Items, error
**/
func (s *LinqFilter) All() (et.Items, error) {
	return s.main.(*Linq).All()
}

/**
* Last
* @param n int
* @return et.Items, error
**/
func (s *LinqFilter) Last(n int) (et.Items, error) {
	return s.main.(*Linq).Last(n)
}

/**
* One
* @return et.Item, error
**/
func (s *LinqFilter) One() (et.Item, error) {
	return s.main.(*Linq).One()
}

/**
* Page
* @param val int
* @return *Linq
**/
func (s *LinqFilter) Page(val int) *Linq {
	return s.main.(*Linq).Page(val)
}

/**
* Rows
* @param val int
* @return et.Items, error
**/
func (s *LinqFilter) Rows(val int) (et.Items, error) {
	return s.main.(*Linq).Rows(val)
}

/**
* List
* @param page, rows int
* @return et.List, error
**/
func (s *LinqFilter) List(page, rows int) (et.List, error) {
	return s.main.(*Linq).List(page, rows)
}

/**
* Query
* @param params et.Json
* @return (et.Items, error)
**/
func (s *LinqFilter) Query(params et.Json) (et.Items, error) {
	return s.main.(*Linq).Query(params)
}
