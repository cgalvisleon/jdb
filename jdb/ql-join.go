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

type QlJoin struct {
	QlWhere
	Ql       *Ql
	TypeJoin TypeJoin
	With     *QlFrom
}

type QlJoins []*QlJoin

/**
* QlJoin
* @param m *Model
* @return *Ql
**/
func (s *Ql) Join(m *Model) *QlJoin {
	with := s.addFrom(m)
	if with == nil {
		return nil
	}

	result := &QlJoin{
		Ql:       s,
		TypeJoin: JoinInner,
		With:     with,
	}

	s.Joins = append(s.Joins, result)

	return result
}

/**
* LeftJoin
* @param m *Model
* @return *Ql
**/
func (s *Ql) LeftJoin(m *Model) *QlJoin {
	result := s.Join(m)
	result.TypeJoin = JoinLeft

	return result
}

/**
* RightJoin
* @param m *Model
* @return *Ql
**/
func (s *Ql) RightJoin(m *Model) *QlJoin {
	result := s.Join(m)
	result.TypeJoin = JoinRight

	return result
}

/**
* FullJoin
* @param m *Model
* @return *Ql
**/
func (s *Ql) FullJoin(m *Model) *QlJoin {
	result := s.Join(m)
	result.TypeJoin = JoinFull

	return result
}

/**
* On
* @param name string
* @return *Ql
**/
func (s *QlJoin) On(val interface{}) *QlJoin {
	switch v := val.(type) {
	case string:
		field := s.Ql.GetField(v)
		if field != nil {
			s.Where(field)
			return s
		}
	}

	s.Where(val)
	return s
}

/**
* And
* @param field string
* @return *QlFilter
**/
func (s *QlJoin) And(val interface{}) *QlJoin {
	switch v := val.(type) {
	case string:
		field := s.Ql.GetField(v)
		if field != nil {
			s.And(field)
			return s
		}
	}

	s.And(val)
	return s
}

/**
* Or
* @param field string
* @return *QlFilter
**/
func (s *QlJoin) Or(val interface{}) *QlJoin {
	switch v := val.(type) {
	case string:
		field := s.Ql.GetField(v)
		if field != nil {
			s.Or(field)
			return s
		}
	}

	s.Or(val)
	return s
}

/**
* Select
* @param fields ...string
* @return *Ql
**/
func (s *QlJoin) Select(fields ...string) *Ql {
	return s.Ql.Select(fields...)
}

/**
* Data
* @param fields ...string
* @return *Ql
**/
func (s *QlJoin) Data(fields ...string) *Ql {
	return s.Ql.Data(fields...)
}

/**
* SetJoins
* @param joins []et.Json
**/
func (s *Ql) setJoins(joins []et.Json) *Ql {
	for _, val := range joins {
		from := val.Str("from")
		model := Jdb.Models[from]
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
func (s *Ql) listJoins() []string {
	result := []string{}
	for _, join := range s.Joins {
		result = append(result, strs.Format(`%s %s AS %s`, join.TypeJoin.Str(), join.With.Table, join.With.As))
		for _, where := range join.Wheres {
			result = append(result, strs.Format(`%s`, where.String()))
		}
	}

	return result
}
