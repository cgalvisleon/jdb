package jdb

import "github.com/cgalvisleon/et/strs"

type Connector int

const (
	Non Connector = iota
	And
	Or
)

func (s Connector) String() string {
	switch s {
	case And:
		return `AND`
	case Or:
		return `OR`
	default:
		return ``
	}
}

type Operator int

const (
	NnN Operator = iota
	Equal
	Neg
	In
	Like
	More
	Less
	MoreEq
	LessEq
	Between
	IsNull
	NotNull
	Search
)

type LinqWhere struct {
	Conector Connector
	A        interface{}
	Operator Operator
	B        []interface{}
}

func (s *LinqWhere) String() string {
	if s.Conector == Non {
		return strs.Format(`%v %v %v`, s.A, s.Operator, s.B)
	} else {
		return strs.Format(`%v %v %v %v`, s.Conector, s.A, s.Operator, s.B)
	}
}

type LinqFilter struct {
	main  interface{}
	where *LinqWhere
}

/**
* NewLinqFilter
* @param main interface{}
* @param val interface{}
* @return *LinqFilter
**/
func NewLinqFilter(main interface{}, val interface{}) *LinqFilter {
	where := &LinqWhere{
		Conector: Non,
		A:        val,
		Operator: NnN,
		B:        make([]interface{}, 0),
	}

	return &LinqFilter{
		main:  main,
		where: where,
	}
}

/**
* Add
* @param val interface{}
* @return interface{}
**/
func (s *LinqFilter) Add(val interface{}) interface{} {
	appendValue := func(linq *Linq, value interface{}) {
		switch v := value.(type) {
		case string:
			field := linq.getSelect(v)
			if field != nil {
				s.where.B = append(s.where.B, field)
			} else {
				s.where.B = append(s.where.B, value)
			}
		default:
			s.where.B = append(s.where.B, value)
		}
	}

	switch m := s.main.(type) {
	case *Linq:
		appendValue(m, val)
		m.Wheres = append(m.Wheres, s.where)
	case *Command:
		field := m.getColumn(val)
		if field != nil {
			s.where.B = append(s.where.B, field)
		} else {
			s.where.B = append(s.where.B, val)
		}
		m.Wheres = append(m.Wheres, s.where)
	case *LinqJoin:
		appendValue(m.Linq, val)
		m.Wheres = append(m.Wheres, s.where)
	}

	return s.main
}

/**
* Eq
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) Eq(val interface{}) *Linq {
	s.where.Operator = Equal
	return s.Add(val).(*Linq)
}

/**
* Neg
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) Neg(val interface{}) *Linq {
	s.where.Operator = Neg
	return s.Add(val).(*Linq)
}

/**
* In
* @param val ...any
* @return *Linq
**/
func (s *LinqFilter) In(val ...any) *Linq {
	s.where.Operator = In
	return s.Add(val).(*Linq)
}

/**
* Like
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) Like(val interface{}) *Linq {
	s.where.Operator = Like
	return s.Add(val).(*Linq)
}

/**
* More
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) More(val interface{}) *Linq {
	s.where.Operator = More
	return s.Add(val).(*Linq)
}

/**
* Less
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) Less(val interface{}) *Linq {
	s.where.Operator = Less
	return s.Add(val).(*Linq)
}

/**
* MoreEq
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) MoreEq(val interface{}) *Linq {
	s.where.Operator = MoreEq
	return s.Add(val).(*Linq)
}

/**
* LessEq
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) LessEs(val interface{}) *Linq {
	s.where.Operator = LessEq
	return s.Add(val).(*Linq)
}

/**
* Search
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) Search(val interface{}) *Linq {
	s.where.Operator = Search
	return s.Add(val).(*Linq)
}

/**
* Between
* @param val1, val2 interface{}
* @return *Linq
**/
func (s *LinqFilter) Between(val1, val2 interface{}) *Linq {
	s.where.Operator = Between
	s.Add(val1)
	return s.Add(val2).(*Linq)
}

/**
* IsNull
* @return *Linq
**/
func (s *LinqFilter) IsNull() *Linq {
	s.where.Operator = IsNull
	return s.main.(*Linq)
}

/**
* NotNull
* @return *Linq
**/
func (s *LinqFilter) NotNull() *Linq {
	s.where.Operator = NotNull
	return s.main.(*Linq)
}

/**
* Where
* @param field string
* @return *Linq
**/
func (s *Linq) Where(field string) *LinqFilter {
	sel := s.getSelect(field)
	if sel != nil {
		return NewLinqFilter(s, sel)
	}

	return NewLinqFilter(s, field)
}

/**
* And
* @param field string
* @return *LinqWheres
**/
func (s *Linq) And(field string) *LinqFilter {
	result := s.Where(field)
	result.where.Conector = And
	return result
}

/**
* And
* @param field string
* @return *LinqWhere
**/
func (s *Linq) Or(field string) *LinqFilter {
	result := s.Where(field)
	result.where.Conector = Or
	return result
}

/**
* ListWheres
* @return []string
**/
func (s *Linq) ListWheres() []string {
	result := []string{}
	for _, val := range s.Wheres {
		result = append(result, val.String())
	}

	return result
}
