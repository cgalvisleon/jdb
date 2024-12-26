package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type TypeJoin int

const (
	JoinInner TypeJoin = iota
	JoinLeft
	JoinRight
	JoinFull
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
	from := s.addFrom(m)
	if from == nil {
		return nil
	}

	result := &LinqJoin{
		Linq:     s,
		TypeJoin: JoinInner,
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
	result.TypeJoin = JoinLeft

	return result
}

/**
* RightJoin
* @param m *Model
* @return *Linq
**/
func (s *Linq) RightJoin(m *Model) *LinqJoin {
	result := s.Join(m)
	result.TypeJoin = JoinRight

	return result
}

/**
* FullJoin
* @param m *Model
* @return *Linq
**/
func (s *Linq) FullJoin(m *Model) *LinqJoin {
	result := s.Join(m)
	result.TypeJoin = JoinFull

	return result
}

/**
* On
* @param col interface{}
* @return *Linq
**/
func (s *LinqJoin) On(name string) *LinqFilter {
	col := s.From.GetField(name)
	if col != nil {
		return NewLinqFilter(s, col)
	}

	return nil
}

/**
* SetJoins
* @param vals []et.Json
**/
func (s *Linq) SetJoins(vals []et.Json) {

}

/**
* ListJoins
* @return []string
**/
func (s *Linq) ListJoins() []string {
	result := []string{}
	for _, join := range s.Joins {
		result = append(result, strs.Format(`%s %s ON %s`, join.TypeJoin, join.From.Table, join.Wheres[0].String()))
	}

	return result
}
