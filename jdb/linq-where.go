package jdb

type Connector int

const (
	Not Connector = iota
	And
	Or
)

type Operator int

const (
	Equal Operator = iota
	Neg
	In
	Like
	More
	Less
	MoreEq
	LessEq
	Between
	IsNull
	Search
)

type LinqWhere struct {
	Conector Connector
	A        interface{}
	Operator Operator
	B        interface{}
}

type LinqFilter struct {
	Linq   *Linq
	Wheres []*LinqWhere
	where  *LinqWhere
}

/**
* NewLinqFilter
* @param l *Linq
* @param wheres []*LinqWhere
* @return *LinqFilter
**/
func NewLinqFilter(l *Linq, whers []*LinqWhere) *LinqFilter {
	return &LinqFilter{
		Linq:   l,
		Wheres: whers,
		where:  &LinqWhere{},
	}
}

func (s *LinqFilter) add(operator Operator, val ...interface{}) *Linq {
	s.where.Operator = operator
	col := s.Linq.getColumn(val)
	if col != nil {
		s.where.B = col
	} else {
		s.where.B = val
	}
	s.Wheres = append(s.Wheres, s.where)

	return s.Linq
}

/**
* Eq
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) Eq(val interface{}) *Linq {
	return s.add(Equal, val)
}

/**
* Neg
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) Neg(val interface{}) *Linq {
	return s.add(Neg, val)
}

/**
* In
* @param val ...interface{}
* @return *Linq
**/
func (s *LinqFilter) In(val ...interface{}) *Linq {
	return s.add(In, val)
}

/**
* Like
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) Like(val interface{}) *Linq {
	return s.add(Like, val)
}

/**
* More
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) More(val interface{}) *Linq {
	return s.add(More, val)
}

/**
* Less
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) Less(val interface{}) *Linq {
	return s.add(Less, val)
}

/**
* MoreEq
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) MoreEq(val interface{}) *Linq {
	return s.add(MoreEq, val)
}

/**
* LessEq
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) LessEs(val interface{}) *Linq {
	return s.add(LessEq, val)
}

/**
* Search
* @param val interface{}
* @return *Linq
**/
func (s *LinqFilter) Search(val interface{}) *Linq {
	return s.add(Search, val)
}

/**
* Between
* @param val1, val2 interface{}
* @return *Linq
**/
func (s *LinqFilter) Between(val1, val2 interface{}) *Linq {
	return s.add(Between, val1, val2)
}

/**
* IsNull
* @return *Linq
**/
func (s *LinqFilter) IsNull() *Linq {
	return s.add(IsNull, nil)
}

/**
* Select
* @param columns ...interface{}
* @return *Linq
**/
func (s *Linq) Where(col interface{}) *LinqFilter {
	where := &LinqWhere{
		Conector: Not,
	}

	_col := s.getColumn(col)
	if _col != nil {
		where.A = _col
	} else {
		where.A = col
	}

	result := &LinqFilter{
		Linq:   s,
		Wheres: s.Wheres,
		where:  where,
	}

	return result
}

/**
* And
* @param col interface{}
* @return *LinqWheres
**/
func (s *Linq) And(col interface{}) *LinqFilter {
	result := s.Where(col)
	result.where.Conector = And
	return result
}

/**
* And
* @param col interface{}
* @return *LinqWhere
**/
func (s *Linq) Or(col interface{}) *LinqFilter {
	result := s.Where(col)
	result.where.Conector = Or
	return result
}
