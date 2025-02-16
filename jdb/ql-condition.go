package jdb

import (
	"fmt"

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
