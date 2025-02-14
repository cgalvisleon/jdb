package jdb

import (
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
	Field     interface{}
	Operator  Operator
	Value     []interface{}
	Language  string
}

/**
* AddValue
* @param val interface{}
* @return *QlCondition
**/
func (s *QlCondition) AddValue(val interface{}) *QlCondition {
	s.Value = append(s.Value, val)

	return s
}

/**
* getValue
* @param val interface{}
* @return string
**/
func (s *QlCondition) GetValue(val interface{}) string {
	switch v := val.(type) {
	case *Field:
		return v.AsField()
	case []interface{}:
		var result string
		for _, w := range v {
			val := s.GetValue(w)
			result = strs.Append(result, strs.Format(`%v`, val), ",")
		}
		return result
	default:
		return strs.Format(`%v`, utility.Quote(v))
	}
}

/**
* GetField
* @return string
**/
func (s *QlCondition) GetField() string {
	return s.GetValue(s.Field)
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
	result = strs.Append(result, s.GetField(), " ")
	result = strs.Append(result, s.Operator.Str(), " ")

	for _, val := range s.Value {
		result = strs.Append(result, s.GetValue(val), " ")
	}

	return result
}

/**
* NewQlCondition
* @params key interface{}
* @return QlWhere
**/
func NewQlCondition(field interface{}) *QlCondition {
	return &QlCondition{
		Connector: NoC,
		Field:     field,
		Operator:  NoP,
		Value:     []interface{}{},
	}
}

type QlWhere struct {
	Wheres   []*QlCondition
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
		Wheres:  []*QlCondition{},
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
	for _, val := range s.Wheres {
		result = strs.Append(result, val.String(), " ")
	}

	return result
}

/**
* Where
* @param val interface{}
* @return *QlWhere
**/
func (s *QlWhere) Where(field interface{}) *QlWhere {
	s.index = len(s.Wheres)
	s.Wheres = append(s.Wheres, NewQlCondition(field))

	return s
}

/**
* And
* @param val interface{}
* @return *QlWhere
**/
func (s *QlWhere) And(val interface{}) *QlWhere {
	s.index = len(s.Wheres)
	where := NewQlCondition(val)
	where.Connector = And
	s.Wheres = append(s.Wheres, where)

	return s
}

/**
* Or
* @param val interface{}
* @return *QlWhere
**/
func (s *QlWhere) Or(val interface{}) *QlWhere {
	s.index = len(s.Wheres)
	where := NewQlCondition(val)
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
	s.Wheres[s.index].Operator = Equal
	s.Wheres[s.index].AddValue(val)

	return s
}

/**
* Neg
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Neg(val interface{}) *QlWhere {
	s.Wheres[s.index].Operator = Neg
	s.Wheres[s.index].AddValue(val)

	return s
}

/**
* In
* @param val ...any
* @return QlWhere
**/
func (s *QlWhere) In(val ...any) *QlWhere {
	s.Wheres[s.index].Operator = In
	s.Wheres[s.index].AddValue(val)

	return s
}

/**
* Like
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Like(val interface{}) *QlWhere {
	s.Wheres[s.index].Operator = Like
	s.Wheres[s.index].AddValue(val)

	return s
}

/**
* More
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) More(val interface{}) *QlWhere {
	s.Wheres[s.index].Operator = More
	s.Wheres[s.index].AddValue(val)

	return s
}

/**
* Less
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Less(val interface{}) *QlWhere {
	s.Wheres[s.index].Operator = Less
	s.Wheres[s.index].AddValue(val)

	return s
}

/**
* MoreEq
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) MoreEq(val interface{}) *QlWhere {
	s.Wheres[s.index].Operator = MoreEq
	s.Wheres[s.index].AddValue(val)

	return s
}

/**
* LessEq
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) LessEq(val interface{}) *QlWhere {
	s.Wheres[s.index].Operator = LessEq
	s.Wheres[s.index].AddValue(val)

	return s
}

/**
* Search
* @param val interface{}
* @return QlWhere
**/
func (s *QlWhere) Full(language string, val interface{}) *QlWhere {
	s.Wheres[s.index].Operator = Search
	s.Wheres[s.index].Language = language
	s.Wheres[s.index].AddValue(val)

	return s
}

/**
* Between
* @param val1, val2 interface{}
* @return QlWhere
**/
func (s *QlWhere) Between(val interface{}) *QlWhere {
	vals, ok := val.([]interface{})
	if !ok {
		return s
	}

	s.Wheres[s.index].Operator = Between
	for _, val := range vals {
		s.Wheres[s.index].AddValue(val)
	}

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
* IsNull
* @return *QlWhere
**/
func (s *QlWhere) IsNull() *QlWhere {
	s.Wheres[s.index].Operator = IsNull
	return s
}

/**
* NotNull
* @return *QlWhere
**/
func (s *QlWhere) NotNull() *QlWhere {
	s.Wheres[s.index].Operator = NotNull
	return s
}

/**
* Language
* @param lan string
* @return *QlWhere
**/
func (s *QlWhere) Language(lan string) *QlWhere {
	s.language = lan
	return s
}

/**
* setWheres
* @param wheres []et.Json
**/
func (s *QlWhere) setWheres(wheres et.Json) *QlWhere {
	setVal := func(val et.Json) {
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
				s.Full(s.language, value)
			}
		}
	}

	and := func(vals []et.Json) {
		for _, val := range vals {
			for key, _ := range val {
				s.And(key)
				setVal(val.Json(key))
			}
		}
	}

	or := func(vals []et.Json) {
		for _, val := range vals {
			for key, _ := range val {
				s.Or(key)
				setVal(val.Json(key))
			}
		}
	}

	where := func(key string, val et.Json) {
		s.Where(key)
		setVal(val)
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
* @return []string
**/
func (s *QlWhere) listWheres() []string {
	result := []string{}
	for _, val := range s.Wheres {
		result = append(result, val.String())
	}

	return result
}
