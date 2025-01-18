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
	*QlFilter
	Ql       *Ql
	TypeJoin TypeJoin
	From     *QlFrom
}

/**
* QlJoin
* @param m *Model
* @return *Ql
**/
func (s *Ql) Join(m *Model) *QlJoin {
	from := s.addFrom(m)
	if from == nil {
		return nil
	}

	result := &QlJoin{
		Ql:       s,
		TypeJoin: JoinInner,
		From:     from,
	}
	result.QlFilter = &QlFilter{
		main:   result,
		Wheres: make([]*QlWhere, 0),
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
func (s *QlJoin) On(field string) *QlFilter {
	col := s.From.GetField(field, false)
	if col != nil {
		s.where = NewQlWhere(col)
	} else {
		s.where = NewQlWhere(field)
	}

	return s.QlFilter
}

/**
* And
* @param field string
* @return *QlFilter
**/
func (s *QlJoin) And(val interface{}) *QlFilter {
	field, ok := val.(string)
	if ok {
		result := s.On(field)
		result.where.Conector = And
	}

	return s.QlFilter
}

/**
* Or
* @param field string
* @return *QlFilter
**/
func (s *QlJoin) Or(val interface{}) *QlFilter {
	field, ok := val.(string)
	if ok {
		result := s.On(field)
		result.where.Conector = Or
	}

	return s.QlFilter
}

/**
* Select
* @param fields ...string
* @return *Ql
**/
func (s *QlJoin) Select(fields ...string) *Ql {
	return s.Ql
}

/**
* Data
* @param fields ...string
* @return *Ql
**/
func (s *QlJoin) Data(fields ...string) *Ql {
	return s.Ql
}

/**
* Exec
* @return et.Items, error
**/
func (s *QlJoin) Exec() (et.Items, error) {
	return et.Items{}, nil
}

/**
* One
* @return et.Item, error
**/
func (s *QlJoin) One() (et.Item, error) {
	return et.Item{}, nil
}

/**
* SetJoins
* @param joins []et.Json
**/
func (s *Ql) setJoins(joins []et.Json) *Ql {
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
func (s *Ql) listJoins() []string {
	result := []string{}
	for _, join := range s.Joins {
		result = append(result, strs.Format(`%s %s AS %s`, join.TypeJoin.Str(), join.From.Table, join.From.As))
		for _, where := range join.Wheres {
			result = append(result, strs.Format(`%s`, where.String()))
		}
	}

	return result
}
