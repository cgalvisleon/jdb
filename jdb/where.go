package jdb

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type Field struct {
	Model     *Model      `json:"model"`
	Field     string      `json:"field"`
	Pattern   int         `json:"pattern"`
	From      string      `json:"from"`
	Name      string      `json:"name"`
	As        string      `json:"as"`
	Type      string      `json:"type"`
	Aggregate string      `json:"aggregate"`
	Value     interface{} `json:"value"`
	Existent  bool        `json:"existent"`
}

func (s *Field) ToJson() et.Json {
	return et.Json{
		"pattern":   s.Pattern,
		"from":      s.From,
		"name":      s.Name,
		"as":        s.As,
		"type":      s.Type,
		"aggregate": s.Aggregate,
		"value":     s.Value,
		"existent":  s.Existent,
	}
}

type Condition struct {
	Field string      `json:"field"`
	Op    string      `json:"op"`
	Value interface{} `json:"value"`
}

/**
* ToJson
* @return et.Json
**/
func (s *Condition) ToJson() et.Json {
	return et.Json{
		s.Field: et.Json{
			s.Op: s.Value,
		},
	}
}

/**
* condition
* @param field string, value interface{}, op string
* @return *Condition
**/
func condition(field string, value interface{}, op string) *Condition {
	return &Condition{
		Field: field,
		Op:    op,
		Value: Quote(value),
	}
}

/**
* Eq
* @param field string, value interface{}
* @return Condition
**/
func Eq(field string, value interface{}) *Condition {
	return condition(field, value, "eq")
}

/**
* Neg
* @param field string, value interface{}
* @return Condition
**/
func Neg(field string, value interface{}) *Condition {
	return condition(field, value, "ne")
}

/**
* Less
* @param field string, value interface{}
* @return Condition
**/
func Less(field string, value interface{}) *Condition {
	return condition(field, value, "less")
}

/**
* LessEq
* @param field string, value interface{}
* @return Condition
**/
func LessEq(field string, value interface{}) *Condition {
	return condition(field, value, "less_eq")
}

/**
* More
* @param field string, value interface{}
* @return Condition
**/
func More(field string, value interface{}) *Condition {
	return condition(field, value, "more")
}

/**
* MoreEq
* @param field string, value interface{}
* @return Condition
**/
func MoreEq(field string, value interface{}) *Condition {
	return condition(field, value, "more_eq")
}

/**
* Like
* @param field string, value interface{}
* @return Condition
**/
func Like(field string, value interface{}) *Condition {
	return condition(field, value, "like")
}

/**
* Ilike
* @param field string, value interface{}
* @return Condition
**/
func Ilike(field string, value interface{}) *Condition {
	return condition(field, value, "ilike")
}

/**
* In
* @param field string, value interface{}
* @return Condition
**/
func In(field string, value interface{}) *Condition {
	return condition(field, value, "in")
}

/**
* NotIn
* @param field string, value interface{}
* @return Condition
**/
func NotIn(field string, value interface{}) *Condition {
	return condition(field, value, "not_in")
}

/**
* Is
* @param field string, value interface{}
* @return Condition
**/
func Is(field string, value interface{}) *Condition {
	return condition(field, value, "is")
}

/**
* IsNot
* @param field string, value interface{}
* @return Condition
**/
func IsNot(field string, value interface{}) *Condition {
	return condition(field, value, "is_not")
}

/**
* Null
* @param field string
* @return Condition
**/
func Null(field string) *Condition {
	return condition(field, nil, "null")
}

/**
* NotNull
* @param field string
* @return Condition
**/
func NotNull(field string) *Condition {
	return condition(field, nil, "not_null")
}

/**
* Between
* @param field string, value []interface{}
* @return Condition
**/
func Between(field string, value []interface{}) *Condition {
	return condition(field, value, "between")
}

/**
* NotBetween
* @param field string, value []interface{}
* @return Condition
**/
func NotBetween(field string, value []interface{}) *Condition {
	return condition(field, value, "not_between")
}

type where struct {
	Froms  map[string]*Model `json:"from"`
	Wheres []et.Json         `json:"where"`
}

/**
* newWhere
* @param model *Model, as string
* @return *where
**/
func newWhere(model *Model, as string) *where {
	result := &where{
		Froms:  map[string]*Model{},
		Wheres: []et.Json{},
	}

	if as == "" {
		as = model.Table
	}

	result.Froms[as] = model
	return result
}

/**
* validField
* @param field string
* @return (*Column, bool)
**/
func (s *where) validField(field *Field) *Field {
	if len(s.Froms) == 0 {
		return field
	}

	if field.From == "" {
		for as, model := range s.Froms {
			col, ok := model.GetColumn(field.Name)
			if ok {
				field.Model = model
				field.From = as
				field.Type = col.Type
				field.Existent = true
				return field
			}
		}
	} else {
		model := s.Froms[field.From]
		col, ok := model.GetColumn(field.Name)
		if ok {
			field.Model = model
			field.Type = col.Type
			field.Existent = true
			return field
		}

		for as, model := range s.Froms {
			if model.Name != field.From {
				continue
			}

			col, ok := model.GetColumn(field.Name)
			if ok {
				field.Model = model
				field.From = as
				field.Type = col.Type
				field.Existent = true
				return field
			}
		}

		for as, model := range s.Froms {
			if model.Table != field.From {
				continue
			}

			col, ok := model.GetColumn(field.Name)
			if ok {
				field.Model = model
				field.From = as
				field.Type = col.Type
				field.Existent = true
				return field
			}
		}
	}

	for as, model := range s.Froms {
		if model.UseAtribs() {
			field.Model = model
			field.From = as
			field.Type = TypeAtribute
			field.Existent = true
			return field
		}
		break
	}

	return field
}

/**
* getField
* @param field string
* @return *Field
**/
func (s *where) getField(field string) *Field {
	result := &Field{
		Field: field,
	}
	pattern1 := regexp.MustCompile(`^([A-Za-z0-9]+)\.([A-Za-z0-9]+):([A-Za-z0-9]+)$`) // from.name:as
	pattern2 := regexp.MustCompile(`^([A-Za-z0-9]+)\.([A-Za-z0-9]+)$`)                // from.name
	pattern3 := regexp.MustCompile(`^([A-Za-z]+)\((.+)\):([A-Za-z0-9]+)$`)            // func(args):as
	pattern4 := regexp.MustCompile(`^([A-Za-z]+)\((.+)\)`)                            // func(args)
	pattern5 := regexp.MustCompile(`^-?\d+(\.\d+)?$`)                                 // number
	if pattern1.MatchString(result.Field) {
		matches := pattern1.FindStringSubmatch(result.Field)
		if len(matches) == 4 {
			result.Pattern = 1
			result.From = matches[1]
			result.Name = matches[2]
			result.As = matches[3]
			result = s.validField(result)
			result.Field = fmt.Sprintf("%s.%s", result.From, result.Name)
		}
	} else if pattern2.MatchString(result.Field) {
		matches := pattern2.FindStringSubmatch(result.Field)
		if len(matches) == 3 {
			result.Pattern = 2
			result.From = matches[1]
			result.Name = matches[2]
			result.As = matches[2]
			result = s.validField(result)
			result.Field = fmt.Sprintf("%s.%s", result.From, result.Name)
		}
	} else if pattern3.MatchString(result.Field) {
		matches := pattern3.FindStringSubmatch(result.Field)
		if len(matches) == 4 {
			result.Pattern = 3
			result.Aggregate = matches[1]
			result.Name = matches[2]
			result.As = matches[3]
			fld := s.getField(result.Name)
			result.From = fld.From
			result.Name = fld.Name
			result = s.validField(result)
			result.Field = fmt.Sprintf("%s(%s)", result.Aggregate, fld.Field)
		}
	} else if pattern4.MatchString(result.Field) {
		matches := pattern4.FindStringSubmatch(result.Field)
		if len(matches) == 3 {
			result.Pattern = 4
			result.Aggregate = matches[1]
			result.Name = matches[2]
			result.As = matches[1]
			fld := s.getField(result.Name)
			result.From = fld.From
			result.Name = fld.Name
			result = s.validField(result)
			result.Field = fmt.Sprintf("%s(%s)", result.Aggregate, fld.Field)
		}
	} else if pattern5.MatchString(result.Field) {
		matches := pattern5.FindStringSubmatch(result.Field)
		if len(matches) == 2 {
			result.Pattern = 5
			result.Value = matches[1]
			result.Name = fmt.Sprintf(`%v`, result.Value)
			result.As = result.Name
			result = s.validField(result)
			result.Field = fmt.Sprintf(`%v`, result.Value)
		}
	} else {
		result.Name = result.Field
		result.As = result.Name
		result = s.validField(result)
		result.Field = strs.Append(result.From, result.Name, ".")
	}

	return result
}

/**
* where
* @param cond Condition
* @return *where
**/
func (s *where) where(cond *Condition, conector string) *where {
	field := s.getField(cond.Field)
	cond.Field = field.Field
	if len(s.Wheres) == 0 {
		s.Wheres = append(s.Wheres, cond.ToJson())
	} else {
		conds := []et.Json{}
		idx := slices.IndexFunc(s.Wheres, func(v et.Json) bool { return strs.Uppcase(v.String(conector)) == strs.Uppcase(conector) })
		if idx != -1 {
			conds = s.Wheres[idx].ArrayJson(conector)
		}

		conds = append(conds, cond.ToJson())
		if idx == -1 {
			s.Wheres = append(s.Wheres, et.Json{
				conector: conds,
			})
		} else {
			s.Wheres[idx][conector] = conds
		}
	}

	return s
}

/**
* Where
* @param cond Condition
* @return *where
**/
func (s *where) Where(cond *Condition) *where {
	return s.where(cond, "and")
}

/**
* And
* @param cond Condition
* @return *where
**/
func (s *where) And(cond *Condition) *where {
	return s.where(cond, "and")
}

/**
* Or
* @param cond Condition
* @return *where
**/
func (s *where) Or(cond *Condition) *where {
	return s.where(cond, "or")
}
