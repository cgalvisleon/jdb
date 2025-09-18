package jdb

import (
	"regexp"
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type Agregation struct {
	Agregation string
	pattern    string
	re         *regexp.Regexp
}

var agregations = map[TypeAgregation]*Agregation{
	Nag:             {Agregation: "", pattern: ""},
	AgregationSum:   {Agregation: "SUM", pattern: `SUM\([a-zA-Z0-9_]+\)$`},
	AgregationCount: {Agregation: "COUNT", pattern: `COUNT\([a-zA-Z0-9_]+\)$`},
	AgregationAvg:   {Agregation: "AVG", pattern: `AVG\([a-zA-Z0-9_]+\)$`},
	AgregationMin:   {Agregation: "MIN", pattern: `MIN\([a-zA-Z0-9_]+\)$`},
	AgregationMax:   {Agregation: "MAX", pattern: `MAX\([a-zA-Z0-9_]+\)$`},
}

/**
* init
**/
func init() {
	for _, agregation := range agregations {
		re, err := regexp.Compile(agregation.pattern)
		if err != nil {
			continue
		}
		agregation.re = re
	}
}

type QlFrom struct {
	*Model
	As string
}

type QlFroms struct {
	Froms []*QlFrom
	index int
}

/**
* newForms
* @return *QlFroms
**/
func newForms() *QlFroms {
	return &QlFroms{
		Froms: make([]*QlFrom, 0),
		index: 65,
	}
}

/**
* setForms
* @param model *Model
* @return *QlFroms
**/
func setForms(model *Model) *QlFroms {
	return &QlFroms{
		Froms: []*QlFrom{
			{
				Model: model,
				As:    "",
			},
		},
		index: 65,
	}
}

/**
* add
* @param m *Model
* @return *QlFrom
**/
func (s *QlFroms) add(m *Model) *QlFrom {
	as := string(rune(s.index))
	from := &QlFrom{
		Model: m,
		As:    as,
	}

	s.Froms = append(s.Froms, from)
	s.index++

	return from
}

/**
* getModel
* @param idx int
* @return *Model
**/
func (s *QlFroms) getModel(idx int) *Model {
	if s.Froms[idx] == nil {
		return nil
	}

	return s.Froms[idx].Model
}

/**
* getForm
* @param idx int
* @return *QlFrom
**/
func (s *QlFroms) getForm(idx int) *QlFrom {
	return s.Froms[idx]
}

/**
* getFormByName
* @param name string
* @return *QlFrom
**/
func (s *QlFroms) getFormByName(name string) *QlFrom {
	for _, from := range s.Froms {
		if from.Name == name {
			return from
		}
	}

	return nil
}

/**
* getField
* @param name string
* @return *Field
**/
func (s *QlFroms) getField(name string, create bool) *Field {
	findField := func(name string) *Field {
		for _, from := range s.Froms {
			field := from.getField(name, create)
			if field != nil {
				field.As = from.As
				return field
			}
		}

		return nil
	}

	for tp, ag := range agregations {
		if ag.re.MatchString(name) {
			n := strs.ReplaceAll(name, []string{ag.Agregation, "(", ")"}, "")
			field := findField(n)
			if field != nil {
				field.Agregation = tp
				return field
			}
		}
	}

	re := regexp.MustCompile(`(?i)\s*AS\s*`)
	list := re.Split(name, -1)
	alias := ""
	if len(list) > 1 {
		name = list[0]
		alias = list[1]
	}

	list = strs.Split(name, ".")
	switch len(list) {
	case 1:
		result := findField(list[0])
		if result != nil && alias != "" {
			result.Alias = alias
		}

		return result
	case 2:
		form := s.getFormByName(list[0])
		if form == nil {
			return nil
		}

		result := form.getField(list[1], create)
		if result != nil && alias != "" {
			result.Alias = alias
		}

		return result
	default:
		return nil
	}
}

/**
* validator
* @param val interface{}
* @return interface{}
**/
func (s *QlFroms) validator(val interface{}) interface{} {
	switch v := val.(type) {
	case string:
		if strings.HasPrefix(v, ":") {
			v = strings.TrimPrefix(v, ":")
			field := s.getField(v, false)
			if field != nil {
				return field
			}
			return nil
		}

		re := regexp.MustCompile(`->>`)
		if re.MatchString(v) {
			filed := newField(v)
			if filed != nil {
				return filed
			}
		}

		field := s.getField(v, false)
		if field != nil {
			return field
		}

		return v
	case *Field:
		return v
	case Field:
		return v
	case *Column:
		return GetField(v)
	case Column:
		return GetField(&v)
	case []interface{}:
		return v
	case []string:
		return v
	case []et.Json:
		return v
	default:
		return v
	}
}
