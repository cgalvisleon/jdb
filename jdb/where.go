package jdb

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

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
		Value: value,
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

type Field struct {
	Pattern     int         `json:"pattern"`
	From        string      `json:"from"`
	Name        string      `json:"name"`
	As          string      `json:"as"`
	Type        string      `json:"type"`
	SourceField string      `json:"source_field"`
	Aggregate   string      `json:"aggregate"`
	Value       interface{} `json:"value"`
}

func (s *Field) ToJson() et.Json {
	return et.Json{
		"pattern":      s.Pattern,
		"from":         s.From,
		"name":         s.Name,
		"as":           s.As,
		"type":         s.Type,
		"source_field": s.SourceField,
		"aggregate":    s.Aggregate,
		"value":        s.Value,
	}
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
* findField
* @param from, name string
* @return (et.Json, bool)
**/
func (s *where) findField(from, name string) (et.Json, bool) {
	if from != "" {
		model := s.Froms[from]
		col, ok := model.GetColumn(name)
		if ok {
			return col, true
		}
	}

	for _, model := range s.Froms {
		if model.Name == from {
			continue
		}

		col, ok := model.GetColumn(name)
		if ok {
			return col, true
		}
	}

	for _, model := range s.Froms {
		if model.Table == from {
			continue
		}

		col, ok := model.GetColumn(name)
		if ok {
			return col, true
		}
	}

	return et.Json{}, false
}

/**
* getField
* @param as string
* @return *Model
**/
func (s *where) getField(field string) Field {
	validField := func(fld Field) Field {
		switch fld.Pattern {
		case 1:
			col, ok := s.findField(fld.From, fld.Name)
			if ok {
				fld.From = col.String("from")

			}

		case 2:
			value.Pattern = 2
		case 3:
			value.Pattern = 3
		case 4:
			value.Pattern = 4
		case 5:
			value.Pattern = 5
		case 6:
			value.Pattern = 6
		}

		return value
	}

	var result Field
	pattern1 := regexp.MustCompile(`^([A-Za-z0-9]+)\.([A-Za-z0-9]+):([A-Za-z0-9]+)$`) // from.name:as
	pattern2 := regexp.MustCompile(`^([A-Za-z0-9]+)\.([A-Za-z0-9]+)$`)                // from.name
	pattern3 := regexp.MustCompile(`^([A-Za-z]+)\((.+)\):([A-Za-z0-9]+)$`)            // func(args):as
	pattern4 := regexp.MustCompile(`^([A-Za-z]+)\((.+)\)`)                            // func(args)
	pattern5 := regexp.MustCompile(`^([A-Za-z]+)$`)                                   // name
	pattern6 := regexp.MustCompile(`^-?\d+(\.\d+)?$`)                                 // number
	if pattern1.MatchString(field) {
		matches := pattern1.FindStringSubmatch(field)
		if len(matches) == 4 {
			result = Field{
				Pattern: 1,
				From:    matches[1],
				Name:    matches[2],
				As:      matches[3],
			}
		}
	} else if pattern2.MatchString(field) {
		matches := pattern2.FindStringSubmatch(field)
		if len(matches) == 3 {
			result = Field{
				Pattern: 2,
				From:    matches[1],
				Name:    matches[2],
				As:      matches[2],
			}
		}
	} else if pattern3.MatchString(field) {
		matches := pattern3.FindStringSubmatch(field)
		if len(matches) == 4 {
			result = Field{
				Pattern:   3,
				As:        matches[3],
				Aggregate: matches[1],
			}
			fld := s.getField(matches[1])
			result.From = fld.From
			result.Name = fld.Name
			result.Type = fld.Type
			result.SourceField = fld.SourceField
		}
	} else if pattern4.MatchString(field) {
		matches := pattern4.FindStringSubmatch(field)
		if len(matches) == 3 {
			result = Field{
				Pattern:   4,
				Aggregate: matches[1],
			}
			fld := s.getField(matches[2])
			result.As = fld.As
			result.From = fld.From
			result.Name = fld.Name
			result.Type = fld.Type
			result.SourceField = fld.SourceField
		}
	} else if pattern5.MatchString(field) {
		matches := pattern5.FindStringSubmatch(field)
		if len(matches) == 2 {
			result = Field{
				Pattern: 5,
				Name:    matches[1],
			}
		}
	} else if pattern6.MatchString(field) {
		matches := pattern6.FindStringSubmatch(field)
		if len(matches) == 2 {
			result = Field{
				Pattern: 6,
				Value:   matches[1],
			}
		}
	} else {
		result = Field{
			Pattern: 5,
			Name:    field,
		}
	}

	result = validField(result)

	return result
}

/**
* where
* @param cond Condition
* @return *where
**/
func (s *where) where(cond *Condition, conector string) *where {
	if s.From != nil {
		f := cond.Field
		col, ok := s.From.GetColumn(f)
		if !ok {
			if s.From.UseAtribs() {
				f = GetFieldName(f)
				cond.Field = fmt.Sprintf("%s:%s", s.From.SourceField, f)
				cond.Field = strs.Append(s.As, cond.Field, ".")
			}
		} else {
			tp := col.String("type")
			if TypeColumn[tp] {
				cond.Field = f
				cond.Field = strs.Append(s.As, cond.Field, ".")
			}
			if TypeAtrib[tp] {
				f = GetFieldName(f)
				cond.Field = fmt.Sprintf("%s:%s", s.From.SourceField, f)
				cond.Field = strs.Append(s.As, cond.Field, ".")
			}
		}
	}

	if len(s.Wheres) == 0 {
		s.Wheres = append(s.Wheres, et.Json{
			cond.Field: et.Json{
				cond.Op: cond.Value,
			},
		})
	} else {
		conds := []et.Json{}
		idx := slices.IndexFunc(s.Wheres, func(v et.Json) bool { return strs.Uppcase(v.String(conector)) == strs.Uppcase(conector) })
		if idx != -1 {
			conds = s.Wheres[idx].ArrayJson(conector)
		}

		conds = append(conds, et.Json{
			cond.Field: et.Json{
				cond.Op: cond.Value,
			},
		})

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
