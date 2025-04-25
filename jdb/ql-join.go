package jdb

import (
	"encoding/json"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
)

type TypeJoin int

const (
	InnerJoin TypeJoin = iota
	LeftJoin
	RightJoin
	FullJoin
)

func (s TypeJoin) Str() string {
	switch s {
	case InnerJoin:
		return "INNER JOIN"
	case LeftJoin:
		return "LEFT JOIN"
	case RightJoin:
		return "RIGHT JOIN"
	case FullJoin:
		return "FULL JOIN"
	}

	return ""
}

type QlJoin struct {
	*QlWhere
	Ql       *Ql      `json:"-"`
	TypeJoin TypeJoin `json:"type_join"`
	With     *QlFrom  `json:"with"`
}

type QlJoins []*QlJoin

/**
* Describe
* @return *et.Json
**/
func (s *QlJoin) Describe() et.Json {
	definition, err := json.Marshal(s)
	if err != nil {
		console.Errorf("QlJoin error: %s", err.Error())
		return et.Json{}
	}

	result := et.Json{}
	err = json.Unmarshal(definition, &result)
	if err != nil {
		console.Errorf("QlJoin error: %s", err.Error())
		return et.Json{}
	}

	result["ql"] = s.Ql.Describe()

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
		field := s.Ql.getField(v)
		if field != nil {
			s.where(field)
			return s
		}
	}

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
		field := s.Ql.getField(v)
		if field != nil {
			s.and(field)
			return s
		}
	}

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
		field := s.Ql.getField(v)
		if field != nil {
			s.or(field)
			return s
		}
	}

	return s
}

/**
* setValue
* @param val et.Json
* @return *QlJoin
**/
func (s *QlJoin) setValue(val et.Json) *QlJoin {
	for key, value := range val {
		switch v := value.(type) {
		case string:
			field := s.Ql.getField(v)
			if field != nil {
				console.Debug("QlJoin:", field.Describe().ToString())
				s.QlWhere.setValue(et.Json{key: field})
				continue
			}

			s.QlWhere.setValue(et.Json{key: v})
		default:
			s.QlWhere.setValue(et.Json{key: v})
		}
	}

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
		QlWhere:  NewQlWhere(),
		Ql:       s,
		TypeJoin: InnerJoin,
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
	result.TypeJoin = LeftJoin

	return result
}

/**
* RightJoin
* @param m *Model
* @return *Ql
**/
func (s *Ql) RightJoin(m *Model) *QlJoin {
	result := s.Join(m)
	result.TypeJoin = RightJoin

	return result
}

/**
* FullJoin
* @param m *Model
* @return *Ql
**/
func (s *Ql) FullJoin(m *Model) *QlJoin {
	result := s.Join(m)
	result.TypeJoin = FullJoin

	return result
}

/**
* SetValue
* @param wheres et.Json
* @return *QlJoin
**/
func (s *QlJoin) setWheres(wheres et.Json) *QlJoin {
	if len(wheres) == 0 {
		return s
	}

	for key := range wheres {
		val := wheres.Json(key)
		s.On(key).setValue(val)
	}

	return s
}

/**
* SetJoins
* @param joins []et.Json
**/
func (s *Ql) setJoins(joins []et.Json) *Ql {
	for _, join := range joins {
		for key := range join {
			with := GetModel(key, false)
			if with != nil {
				val := join.Json(key)
				s.Join(with).setWheres(val)
			}
		}
	}

	return s
}

/**
* listJoins
* @return []et.Json
**/
func (s *Ql) listJoins() []et.Json {
	result := []et.Json{}
	for _, join := range s.Joins {
		item := et.Json{
			join.With.Name: join.listWheres(s.asField),
		}
		result = append(result, item)
	}

	return result
}
