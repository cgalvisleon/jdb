package jdb

type TypeJoin int

const (
	TypeJoinInner TypeJoin = iota
	TypeJoinLeft
	TypeJoinRight
	TypeJoinFull
)

type LinqJoin struct {
	Linq     *Linq
	TypeJoin TypeJoin
	From     *LinqFrom
	Wheres   []*LinqWhere
}

/**
* LinqJoin
* @param m *Model
* @return *Linq
**/
func (s *Linq) Join(m *Model) *LinqJoin {
	from := s.addFrom(*m)
	if from == nil {
		return nil
	}

	result := &LinqJoin{
		Linq:     s,
		TypeJoin: TypeJoinInner,
		From:     from,
		Wheres:   make([]*LinqWhere, 0),
	}
	s.Joins = append(s.Joins, result)

	return result
}

/**
* LeftJoin
* @param m *Model
* @return *Linq
**/
func (s *Linq) LeftJoin(m *Model) *LinqJoin {
	result := s.Join(m)
	result.TypeJoin = TypeJoinLeft

	return result
}

/**
* RightJoin
* @param m *Model
* @return *Linq
**/
func (s *Linq) RightJoin(m *Model) *LinqJoin {
	result := s.Join(m)
	result.TypeJoin = TypeJoinRight

	return result
}

/**
* FullJoin
* @param m *Model
* @return *Linq
**/
func (s *Linq) FullJoin(m *Model) *LinqJoin {
	result := s.Join(m)
	result.TypeJoin = TypeJoinFull

	return result
}

/**
* On
* @param col interface{}
* @return *Linq
**/
func (s *LinqJoin) On(col interface{}) *LinqFilter {
	result := &LinqFilter{
		Linq:   s.Linq,
		Wheres: s.Wheres,
		where:  &LinqWhere{},
	}

	return result
}
