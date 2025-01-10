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
	*LinqFilter
	Linq     *Linq
	TypeJoin TypeJoin
	From     *LinqFrom
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
	}
	result.LinqFilter = &LinqFilter{
		main:   result,
		Wheres: make([]*LinqWhere, 0),
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
	col := s.From.GetField(field, false)
	if col != nil {
		s.where = NewLinqWhere(col)
	} else {
		s.where = NewLinqWhere(field)
	}

	return s.LinqFilter
}

/**
* And
* @param field string
* @return *LinqFilter
**/
func (s *LinqJoin) And(val interface{}) *LinqFilter {
	field, ok := val.(string)
	if ok {
		result := s.On(field)
		result.where.Conector = And
	}

	return s.LinqFilter
}

/**
* Or
* @param field string
* @return *LinqFilter
**/
func (s *LinqJoin) Or(val interface{}) *LinqFilter {
	field, ok := val.(string)
	if ok {
		result := s.On(field)
		result.where.Conector = Or
	}

	return s.LinqFilter
}

/**
* Select
* @param fields ...string
* @return *Linq
**/
func (s *LinqJoin) Select(fields ...string) *Linq {
	return s.Linq
}

/**
* Data
* @param fields ...string
* @return *Linq
**/
func (s *LinqJoin) Data(fields ...string) *Linq {
	return s.Linq
}

/**
* Exec
* @return et.Items, error
**/
func (s *LinqJoin) Exec() (et.Items, error) {
	return et.Items{}, nil
}

/**
* One
* @return et.Item, error
**/
func (s *LinqJoin) One() (et.Item, error) {
	return et.Item{}, nil
}

/**
* SetJoins
* @param joins []et.Json
**/
func (s *Linq) setJoins(joins []et.Json) *Linq {
	for _, val := range joins {
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
