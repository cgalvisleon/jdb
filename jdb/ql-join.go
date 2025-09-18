package jdb

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

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
* Serialize
* @return []byte, error
**/
func (s *QlJoin) Serialize() ([]byte, error) {
	result, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}

	return result, nil
}

/**
* Describe
* @return *et.Json
**/
func (s *QlJoin) Describe() et.Json {
	definition, err := s.Serialize()
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
* @param val string
* @return *QlJoin
**/
func (s *QlJoin) On(val string) *QlJoin {
	field := s.Ql.getField(val)
	if field != nil {
		s.setWhere(field)
	}

	return s
}

/**
* And
* @param val interface{}
* @return *QlJoin
**/
func (s *QlJoin) And(val interface{}) *QlJoin {
	val = s.Ql.validator(val)
	if val != nil {
		s.setAnd(val)
	}

	return s
}

/**
* Or
* @param val interface{}
* @return *QlJoin
**/
func (s *QlJoin) Or(val interface{}) *QlJoin {
	val = s.Ql.validator(val)
	if val != nil {
		s.setOr(val)
	}

	return s
}

/**
* Select
* @param fields ...interface{}
* @return *Ql
**/
func (s *QlJoin) Select(fields ...interface{}) *Ql {
	return s.Ql.Select(fields...)
}

/**
* Data
* @param fields ...interface{}
* @return *Ql
**/
func (s *QlJoin) Data(fields ...interface{}) *Ql {
	return s.Ql.Data(fields...)
}

/**
* QlJoin
* @param name interface{}
* @return *Ql
**/
func (s *Ql) Join(name interface{}) *QlJoin {
	var model *Model
	switch v := name.(type) {
	case *Model:
		model = v
	default:
		str := fmt.Sprintf("%v", v)
		model = s.Db.GetModel(str)
	}

	with := s.Froms.add(model)
	result := &QlJoin{
		QlWhere:  newQlWhere(s.validator),
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
* setWheres
* @param wheres et.Json
* @return *QlJoin
**/
func (s *QlJoin) setWheres(wheres et.Json) *QlJoin {
	if len(wheres) == 0 {
		return s
	}

	and := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				s.setAnd(key)
				s.setValue(val.Json(key))
			}
		}
	}

	or := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				s.setOr(key)
				s.setValue(val.Json(key))
			}
		}
	}

	for key := range wheres {
		key = strings.ToLower(key)
		if slices.Contains([]string{"and", "or"}, key) {
			continue
		}

		val := wheres.Json(key)
		s.On(key).setValue(val)
	}

	for key := range wheres {
		switch strings.ToLower(key) {
		case "and":
			vals := wheres.ArrayJson(key)
			and(vals)
		case "or":
			vals := wheres.ArrayJson(key)
			or(vals)
		}
	}

	return s
}

/**
* SetJoins
* @param joins []et.Json
**/
func (s *Ql) SetJoins(joins []et.Json) *Ql {
	for _, join := range joins {
		for key := range join {
			with, err := LoadModel(s.Db, key)
			if err != nil {
				console.Errorf("Ql error: %s", err.Error())
				continue
			}

			if with != nil {
				val := join.Json(key)
				s.Join(with).setWheres(val)
			}
		}
	}

	return s
}

/**
* getJoins
* @return []et.Json
**/
func (s *Ql) getJoins() []et.Json {
	result := []et.Json{}
	for _, join := range s.Joins {
		item := et.Json{
			join.With.Name: join.getWheres(),
		}
		result = append(result, item)
	}

	return result
}
