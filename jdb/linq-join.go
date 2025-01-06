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

func (s TypeJoin) Str() string {
	switch s {
	case JoinInner:
		return "INNER JOIN"
	case JoinLeft:
		return "LEFT JOIN"
	case JoinRight:
		return "RIGHT JOIN"
	case JoinFull:
		return "FULL JOIN"
	}

	return ""
}

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
* @param name string
* @return *Linq
**/
func (s *LinqJoin) On(field string) *LinqFilter {
	col := s.From.GetField(field)
	if col != nil {
		return NewLinqFilter(s, col)
	}

	return NewLinqFilter(s, field)
}

/**
* And
* @param val interface{}
* @return *LinqFilter
**/
func (s *LinqJoin) And(val interface{}) *LinqFilter {
	field, ok := val.(string)
	if !ok {
		return nil
	}

	result := s.On(field)
	result.where.Conector = And

	return result
}

/**
* Or
* @param field string
* @return *LinqFilter
**/
func (s *LinqJoin) Or(val interface{}) *LinqFilter {
	field, ok := val.(string)
	if !ok {
		return nil
	}

	result := s.On(field)
	result.where.Conector = And

	return result
}

/**
* Select
* @param fields ...string
* @return FilterTo
**/
func (s *LinqJoin) Select(fields ...string) *Linq {
	return s.Linq
}

/**
* Data
* @param fields ...string
* @return FilterTo
**/
func (s *LinqJoin) Data(fields ...string) *Linq {
	return s.Linq
}

/**
* Return
* @param fields ...string
* @return *Command
**/
func (s *LinqJoin) Return(fields ...string) *Command {
	return nil
}

/**
* SetJoins
* @param vals []et.Json
**/
func (s *Linq) setJoins(vals []et.Json) *Linq {
	for _, val := range vals {
		from := val.Str("from")
		model := models[from]
		if model != nil {
			on := val.Json("on")
			key := strs.Format(`%s.%s`, from, on.Str("key"))
			to := on.Str("to")
			foreign := on.Str("foreignKey")
			foreignKey := strs.Format(`%s.%s`, to, foreign)
			s.Join(model).On(key).
				Eq(foreignKey)
		}
	}

	return s
}

/**
* listJoins
* @return []string
**/
func (s *Linq) listJoins() []string {
	result := []string{}
	for _, join := range s.Joins {
		result = append(result, strs.Format(`%s %s AS %s`, join.TypeJoin.Str(), join.From.Table, join.From.As))
		for _, where := range join.Wheres {
			result = append(result, strs.Format(`%s`, where.String()))
		}
	}

	return result
}
