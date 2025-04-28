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
	case IsNull:
		return "is null"
	case NotNull:
		return "is not null"
	case Search:
		return "search"
	default:
		return "Any"
	}
}

/**
* Name
* @return string
**/
func (s *Operator) Name() string {
	switch *s {
	case Equal:
		return "Equal"
	case Neg:
		return "Neg"
	case In:
		return "In"
	case Like:
		return "Like"
	case More:
		return "More"
	case Less:
		return "Less"
	case MoreEq:
		return "MoreEq"
	case LessEq:
		return "LessEq"
	case Between:
		return "Between"
	case IsNull:
		return "IsNull"
	case NotNull:
		return "NotNull"
	case Search:
		return "Search"
	default:
		return "Any"
	}
}

/**
* Command
* @return string
**/
func (s *Operator) Command() string {
	switch *s {
	case Equal:
		return "eq"
	case Neg:
		return "neg"
	case In:
		return "in"
	case Like:
		return "like"
	case More:
		return "more"
	case Less:
		return "less"
	case MoreEq:
		return "moreEq"
	case LessEq:
		return "lessEq"
	case Between:
		return "between"
	case IsNull:
		return "isNull"
	case NotNull:
		return "notNull"
	case Search:
		return "search"
	default:
		return "any"
	}
}

/**
* OperatorToCommand
* @param op Operator
* @return string
**/
func OperatorToCommand(op Operator) string {
	return op.Command()
}

/**
* StrToOperator
* @param str string
* @return Operator
**/
func StrToOperator(str string) Operator {
	switch str {
	case "eq":
		return Equal
	case "neg":
		return Neg
	case "in":
		return In
	case "like":
		return Like
	case "more":
		return More
	case "less":
		return Less
	case "moreEq":
		return MoreEq
	case "lessEq":
		return LessEq
	case "between":
		return Between
	case "isNull":
		return IsNull
	case "notNull":
		return NotNull
	case "search":
		return Search
	default:
		return NoP
	}
}

type QlCondition struct {
	Connector Connector
	Field     interface{}
	Operator  Operator
	Value     interface{}
	Language  string
}

/**
* NewQlCondition
* @params field interface{}
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

/**
* setVal
* @param val interface{}
* @return *QlCondition
**/
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
* getField
* @return string
**/
func (s *QlCondition) getField() string {
	switch v := s.Value.(type) {
	case *Field:
		return v.asName()
	case Field:
		return v.asName()
	default:
		return fmt.Sprintf(`%v`, utility.Quote(v))
	}
}

/**
* getValue
* @return string
**/
func (s *QlCondition) getValue() string {
	switch v := s.Value.(type) {
	case *Field:
		return v.asName()
	case Field:
		return v.asName()
	default:
		return fmt.Sprintf(`%v`, utility.Quote(v))
	}
}

/**
* String
* @return string
**/
func (s *QlCondition) string() string {
	var result string

	if s.Connector != NoC {
		result = strs.Append(result, s.Connector.Str(), " ")
	}
	result = strs.Append(result, s.getField(), " ")
	result = strs.Append(result, s.Operator.Str(), " ")
	result = strs.Append(result, s.getValue(), " ")

	return result
}
