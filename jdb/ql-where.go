package jdb

import (
	"slices"
	"strings"

	"github.com/cgalvisleon/et/et"
)

/**
* Where
* @param val interface{}
* @return *Ql
**/
func (s *Ql) Where(val interface{}) *Ql {
	s.QlWhere.Where(val)

	return s
}

/**
* And
* @param val interface{}
* @return *Ql
**/
func (s *Ql) And(val interface{}) *Ql {
	s.QlWhere.And(val)

	return s
}

/**
* Or
* @param val interface{}
* @return *Ql
**/
func (s *Ql) Or(val interface{}) *Ql {
	s.QlWhere.Or(val)

	return s
}

/**
* Eq
* @param val interface{}
* @return *Ql
**/
func (s *Ql) Eq(val interface{}) *Ql {
	s.QlWhere.Eq(val)

	return s
}

/**
* Neg
* @param val interface{}
* @return *Ql
**/
func (s *Ql) Neg(val interface{}) *Ql {
	s.QlWhere.Neg(val)

	return s
}

/**
* In
* @param val ...any
* @return *Ql
**/
func (s *Ql) In(val ...any) *Ql {
	s.QlWhere.In(val...)

	return s
}

/**
* Like
* @param val interface{}
* @return *Ql
**/
func (s *Ql) Like(val interface{}) *Ql {
	s.QlWhere.Like(val)

	return s
}

/**
* More
* @param val interface{}
* @return *Ql
**/
func (s *Ql) More(val interface{}) *Ql {
	s.QlWhere.More(val)

	return s
}

/**
* Less
* @param val interface{}
* @return *Ql
**/
func (s *Ql) Less(val interface{}) *Ql {
	s.QlWhere.Less(val)

	return s
}

/**
* MoreEq
* @param val interface{}
* @return *Ql
**/
func (s *Ql) MoreEq(val interface{}) *Ql {
	s.QlWhere.MoreEq(val)

	return s
}

/**
* LessEq
* @param val interface{}
* @return *Ql
**/
func (s *Ql) LessEq(val interface{}) *Ql {
	s.QlWhere.LessEq(val)

	return s
}

/*
*
* Between
* @param vals interface{}
* @return *Ql
**/
func (s *Ql) Between(vals interface{}) *Ql {
	s.QlWhere.Between(vals)

	return s
}

/**
* Search
* @param language string, val interface{}
* @return *Ql
**/
func (s *Ql) Search(language string, val interface{}) *Ql {
	s.QlWhere.Search(language, val)

	return s
}

/**
* IsNull
* @return *Ql
**/
func (s *Ql) IsNull() *Ql {
	s.QlWhere.IsNull()

	return s
}

/**
* NotNull
* @return *Ql
**/
func (s *Ql) NotNull() *Ql {
	s.QlWhere.NotNull()

	return s
}

/**
* Debug
* @param v bool
* @return *Ql
**/
func (s *Ql) Debug() *Ql {
	s.QlWhere.Debug()

	return s
}

/**
* setDebug
* @param debug bool
* @return *Ql
**/
func (s *Ql) SetDebug(debug bool) *Ql {
	s.QlWhere.setDebug(debug)

	return s
}

/**
* SetWheres
* @param wheres et.Json
* @return *Ql
**/
func (s *Ql) SetWheres(wheres et.Json) *Ql {
	and := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				s.And(key).setValue(val.Json(key))
			}
		}
	}

	or := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				s.Or(key).setValue(val.Json(key))
			}
		}
	}

	for key := range wheres {
		key = strings.ToLower(key)
		if slices.Contains([]string{"and", "or"}, key) {
			continue
		}

		val := wheres.Json(key)
		s.Where(key).setValue(val)
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
