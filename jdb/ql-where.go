package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type QlWhere struct {
	Wheres   []*QlCondition
	history  bool
	language string
	IsDebug  bool
}

/**
* NewQlWhere
* @return *QlWhere
**/
func NewQlWhere() *QlWhere {
	return &QlWhere{
		Wheres:  []*QlCondition{},
		history: false,
		IsDebug: false,
	}
}

func (s *QlWhere) whr() *QlCondition {
	idx := len(s.Wheres)
	if idx <= 0 {
		return nil
	}

	return s.Wheres[idx-1]
}

/**
* String
* @return string
**/
func (s *QlWhere) String() string {
	var result string
	for _, val := range s.Wheres {
		result = strs.Append(result, val.String(), " ")
	}

	return result
}

/**
* where
* @param val field *Field
* @return *QlWhere
**/
func (s *QlWhere) where(field *Field) *QlWhere {
	if field == nil {
		return s
	}

	s.Wheres = append(s.Wheres, NewQlCondition(field))

	return s
}

/**
* and
* @param val field *Field
* @return *QlWhere
**/
func (s *QlWhere) and(field *Field) *QlWhere {
	where := NewQlCondition(field)
	where.Connector = And
	s.Wheres = append(s.Wheres, where)

	return s
}

/**
* or
* @param val field *Field
* @return *QlWhere
**/
func (s *QlWhere) or(field *Field) *QlWhere {
	where := NewQlCondition(field)
	where.Connector = Or
	s.Wheres = append(s.Wheres, where)

	return s
}

/**
* Eq
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Eq(val interface{}) *QlWhere {
	whr := s.whr()
	if whr == nil {
		return s
	}

	whr.Operator = Equal
	whr.setVal(val)

	return s
}

/**
* Neg
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Neg(val interface{}) *QlWhere {
	whr := s.whr()
	if whr == nil {
		return s
	}

	whr.Operator = Neg
	whr.setVal(val)

	return s
}

/**
* In
* @param val ...any
* @return QlWhere
**/
func (s *QlWhere) In(val ...any) *QlWhere {
	whr := s.whr()
	if whr == nil {
		return s
	}

	whr.Operator = In
	whr.setVal(val)

	return s
}

/**
* Like
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Like(val interface{}) *QlWhere {
	whr := s.whr()
	if whr == nil {
		return s
	}

	whr.Operator = Like
	whr.setVal(val)

	return s
}

/**
* More
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) More(val interface{}) *QlWhere {
	whr := s.whr()
	if whr == nil {
		return s
	}

	whr.Operator = More
	whr.setVal(val)

	return s
}

/**
* Less
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Less(val interface{}) *QlWhere {
	whr := s.whr()
	if whr == nil {
		return s
	}

	whr.Operator = Less
	whr.setVal(val)

	return s
}

/**
* MoreEq
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) MoreEq(val interface{}) *QlWhere {
	whr := s.whr()
	if whr == nil {
		return s
	}

	whr.Operator = MoreEq
	whr.setVal(val)

	return s
}

/**
* LessEq
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) LessEq(val interface{}) *QlWhere {
	whr := s.whr()
	if whr == nil {
		return s
	}

	whr.Operator = LessEq
	whr.setVal(val)

	return s
}

/**
* Between
* @param val1, val2 interface{}
* @return QlWhere
**/
func (s *QlWhere) Between(vals interface{}) *QlWhere {
	val, ok := vals.([]interface{})
	if !ok {
		return s
	}

	whr := s.whr()
	if whr == nil {
		return s
	}

	whr.Operator = Between
	whr.setVal(val)

	return s
}

/**
* Search
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Search(language string, val interface{}) *QlWhere {
	whr := s.whr()
	if whr == nil {
		return s
	}

	whr.Operator = Search
	whr.Language = language
	whr.setVal(val)

	return s
}

/**
* IsNull
* @return *QlWhere
**/
func (s *QlWhere) IsNull() *QlWhere {
	whr := s.whr()
	if whr == nil {
		return s
	}

	whr.Operator = IsNull

	return s
}

/**
* NotNull
* @return *QlWhere
**/
func (s *QlWhere) NotNull() *QlWhere {
	whr := s.whr()
	if whr == nil {
		return s
	}

	whr.Operator = NotNull

	return s
}

/**
* History
* @param v bool
* @return *Command
**/
func (s *QlWhere) History(v bool) *QlWhere {
	s.history = v

	return s
}

/**
* Debug
* @param v bool
* @return *Command
**/
func (s *QlWhere) Debug() *QlWhere {
	s.IsDebug = true

	return s
}

/**
* setValue
* @param val et.Json
* @return *QlWhere
**/
func (s *QlWhere) setValue(val et.Json) *QlWhere {
	for key, value := range val {
		switch key {
		case "eq":
			s.Eq(value)
		case "neg":
			s.Neg(value)
		case "in":
			s.In(value)
		case "like":
			s.Like(value)
		case "more":
			s.More(value)
		case "less":
			s.Less(value)
		case "moreEq":
			s.MoreEq(value)
		case "lessEq":
			s.LessEq(value)
		case "between":
			s.Between(value)
		case "isNull":
			s.IsNull()
		case "notNull":
			s.NotNull()
		case "search":
			s.Search(s.language, value)
		}
	}

	return s
}

/**
* setWheres
* @param wheres et.Json, findField func(name string) *Field
* @return *QlWhere
**/
func (s *QlWhere) setWheres(wheres et.Json, findField func(name string) *Field) *QlWhere {
	if len(wheres) == 0 {
		return s
	}

	and := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				field := findField(key)
				if field != nil {
					s.and(field)
					s.setValue(val.Json(key))
				}
			}
		}
	}

	or := func(vals []et.Json) {
		for _, val := range vals {
			for key := range val {
				field := findField(key)
				if field != nil {
					s.or(field)
					s.setValue(val.Json(key))
				}
			}
		}
	}

	where := func(key string, val et.Json) {
		field := findField(key)
		if field != nil {
			s.where(field)
			s.setValue(val)
		}
	}

	for key := range wheres {
		if slices.Contains([]string{"and", "AND", "or", "OR"}, key) {
			continue
		}

		val := wheres.Json(key)
		where(key, val)
	}

	for key := range wheres {
		switch key {
		case "and", "AND":
			vals := wheres.ArrayJson(key)
			and(vals)
		case "or", "OR":
			vals := wheres.ArrayJson(key)
			or(vals)
		}
	}

	return s
}

/**
* listWheres
* @param asField func(field *Field) string
* @return et.Json
**/
func (s *QlWhere) listWheres(asField func(field *Field) string) et.Json {
	result := et.Json{}
	and := []et.Json{}
	or := []et.Json{}
	for i, con := range s.Wheres {
		if con.Field == nil {
			continue
		}

		field := asField(con.Field)
		def := et.Json{con.Operator.Str(): con.ValStr()}
		if con.Connector == And {
			and = append(and, et.Json{field: def})
		} else if con.Connector == Or {
			or = append(or, et.Json{field: def})
		} else if i == 0 {
			result.Set(field, def)
		}
	}

	if len(and) > 0 {
		result.Set("and", and)
	}
	if len(or) > 0 {
		result.Set("or", or)
	}

	return result
}

/**
* Where
* @param val interface{}
* @return *Ql
**/
func (s *Ql) Where(val string) *Ql {
	field := s.getField(val)
	if field != nil {
		s.where(field)
	}

	return s
}

/**
* And
* @param val interface{}
* @return *Ql
**/
func (s *Ql) And(val string) *Ql {
	field := s.getField(val)
	if field != nil {
		s.and(field)
	}

	return s
}

/**
* Or
* @param val interface{}
* @return *Ql
**/
func (s *Ql) Or(val string) *Ql {
	field := s.getField(val)
	if field != nil {
		s.or(field)
	}

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
* History
* @param v bool
* @return *Ql
**/
func (s *Ql) History(v bool) *Ql {
	s.QlWhere.History(v)

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
* setWheres
* @param wheres []et.Json
*
 */
func (s *Ql) setWheres(wheres et.Json) *Ql {
	s.QlWhere.setWheres(wheres, s.getField)

	return s
}

/**
* listWheres
* @return et.Json
**/
func (s *Ql) listWheres() et.Json {
	return s.QlWhere.listWheres(s.asField)
}
