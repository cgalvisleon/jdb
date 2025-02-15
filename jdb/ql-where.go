package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
)

type Connector int

const (
	NoC Connector = iota
	And
	Or
)

func (s Connector) Str() string {
	switch s {
	case And:
		return "and"
	case Or:
		return "or"
	default:
		return ""
	}
}

type Operator int

const (
	NoP Operator = iota
	Equal
	Neg
	In
	Like
	More
	Less
	MoreEq
	LessEq
	Between
	IsNull
	NotNull
	Search
)

func (s Operator) Str() string {
	switch s {
	case Equal:
		return "="
	case Neg:
		return "!="
	case In:
		return "in"
	case Like:
		return "like"
	case More:
		return ">"
	case Less:
		return "<"
	case MoreEq:
		return ">="
	case LessEq:
		return "<="
	case Between:
		return "between"
	case
		IsNull:
		return "is null"
	case NotNull:
		return "is not null"
	case Search:
		return "search"
	default:
		return ""
	}
}

type QlCondition struct {
	Connector Connector
	Field     *Field
	Operator  Operator
	Value     interface{}
	Language  string
}

func (s *QlCondition) setVal(val interface{}) {
	switch v := val.(type) {
	case *Field:
		s.Value = v
	case Field:
		s.Value = v
	case Column:
		s.Value = v.GetField()
	case *Column:
		s.Value = v.GetField()
	default:
		s.Value = val
	}
}

/**
* GetValue
* @param val interface{}
* @return string
**/
func (s *QlCondition) GetValue() interface{} {
	switch v := s.Value.(type) {
	case *Field:
		return v.AsName()
	case Field:
		return v.AsName()
	default:
		return strs.Format(`%v`, utility.Quote(v))
	}
}

/**
* ValStr
* @return *string
**/
func (s *QlCondition) ValStr() string {
	return fmt.Sprintf(`%v`, s.GetValue())
}

/**
* String
* @return string
**/
func (s *QlCondition) String() string {
	var result string

	if s.Connector != NoC {
		result = strs.Append(result, s.Connector.Str(), " ")
	}
	result = strs.Append(result, s.Field.AsName(), " ")
	result = strs.Append(result, s.Operator.Str(), " ")
	result = strs.Append(result, s.ValStr(), " ")

	return result
}

/**
* NewQlCondition
* @params key interface{}
* @return QlWhere
**/
func NewQlCondition(field *Field) *QlCondition {
	return &QlCondition{
		Connector: NoC,
		Field:     field,
		Operator:  NoP,
		Value:     []interface{}{},
	}
}

type QlWhere struct {
	wheres   []*QlCondition
	index    int
	history  bool
	debug    bool
	language string
}

/**
* NewQlWhere
* @return *QlWhere
**/
func NewQlWhere() *QlWhere {
	return &QlWhere{
		wheres:  []*QlCondition{},
		index:   0,
		history: false,
		debug:   false,
	}
}

/**
* String
* @return string
**/
func (s *QlWhere) String() string {
	var result string
	for _, val := range s.wheres {
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
	s.index = len(s.wheres)
	s.wheres = append(s.wheres, NewQlCondition(field))

	return s
}

/**
* and
* @param val field *Field
* @return *QlWhere
**/
func (s *QlWhere) and(field *Field) *QlWhere {
	s.index = len(s.wheres)
	where := NewQlCondition(field)
	where.Connector = And
	s.wheres = append(s.wheres, where)

	return s
}

/**
* or
* @param val field *Field
* @return *QlWhere
**/
func (s *QlWhere) or(field *Field) *QlWhere {
	s.index = len(s.wheres)
	where := NewQlCondition(field)
	where.Connector = Or
	s.wheres = append(s.wheres, where)

	return s
}

/**
* Eq
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Eq(val interface{}) *QlWhere {
	s.wheres[s.index].Operator = Equal
	s.wheres[s.index].setVal(val)

	return s
}

/**
* Neg
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Neg(val interface{}) *QlWhere {
	s.wheres[s.index].Operator = Neg
	s.wheres[s.index].setVal(val)

	return s
}

/**
* In
* @param val ...any
* @return QlWhere
**/
func (s *QlWhere) In(val ...any) *QlWhere {
	s.wheres[s.index].Operator = In
	s.wheres[s.index].setVal(val)

	return s
}

/**
* Like
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Like(val interface{}) *QlWhere {
	s.wheres[s.index].Operator = Like
	s.wheres[s.index].setVal(val)

	return s
}

/**
* More
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) More(val interface{}) *QlWhere {
	s.wheres[s.index].Operator = More
	s.wheres[s.index].setVal(val)

	return s
}

/**
* Less
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Less(val interface{}) *QlWhere {
	s.wheres[s.index].Operator = Less
	s.wheres[s.index].setVal(val)

	return s
}

/**
* MoreEq
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) MoreEq(val interface{}) *QlWhere {
	s.wheres[s.index].Operator = MoreEq
	s.wheres[s.index].setVal(val)

	return s
}

/**
* LessEq
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) LessEq(val interface{}) *QlWhere {
	s.wheres[s.index].Operator = LessEq
	s.wheres[s.index].setVal(val)

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

	s.wheres[s.index].Operator = Between
	s.wheres[s.index].setVal(val)

	return s
}

/**
* Search
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Search(language string, val interface{}) *QlWhere {
	s.wheres[s.index].Operator = Search
	s.wheres[s.index].Language = language
	s.wheres[s.index].setVal(val)

	return s
}

/**
* IsNull
* @return *QlWhere
**/
func (s *QlWhere) IsNull() *QlWhere {
	s.wheres[s.index].Operator = IsNull
	return s
}

/**
* NotNull
* @return *QlWhere
**/
func (s *QlWhere) NotNull() *QlWhere {
	s.wheres[s.index].Operator = NotNull
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
	s.debug = true

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
* @param wheres et.Json
**/
func (s *QlWhere) setWheres(wheres et.Json, findField func(name string) *Field) *QlWhere {
	and := func(vals []et.Json) {
		for _, val := range vals {
			for key, _ := range val {
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
			for key, _ := range val {
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

	for key, _ := range wheres {
		if key == "and" {
			vals := wheres.ArrayJson(key)
			and(vals)
		} else if key == "or" {
			vals := wheres.ArrayJson(key)
			or(vals)
		} else {
			val := wheres.Json(key)
			where(key, val)
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
	for i, con := range s.wheres {
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
