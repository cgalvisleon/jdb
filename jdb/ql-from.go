package jdb

import (
	"fmt"
	"sync"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

func helpQl() et.Json {
	return et.Json{
		"help": "Las consultas en jdb, se realizan siguiendo el siguiente formato:",
		"params": et.Json{
			"from": "model",
			"data": []interface{}{
				"name",
				"status_id",
				"kinds.name:status",
				et.Json{
					"folders": et.Json{
						"select": []interface{}{
							"name",
						},
						"page": 1,
						"rows": 30,
						"list": true,
					},
				},
			},
			"join": []et.Json{
				{
					"kinds": et.Json{
						"status_id": et.Json{
							"eq": "kinds.id",
						},
					},
					"AND": []et.Json{},
					"OR":  []et.Json{},
				},
			},
			"where": et.Json{
				"status_id": et.Json{
					"eq": "kinds.id",
				},
				"AND": []et.Json{
					{
						"name": et.Json{
							"eq": "v:name",
						},
					},
				},
				"OR": []et.Json{},
			},
			"group_by": []string{"name"},
			"having": et.Json{
				"name": et.Json{
					"eq": "name",
				},
				"AND": []et.Json{},
				"OR":  []et.Json{},
			},
			"order_by": et.Json{
				"ASC":  []string{"name"},
				"DESC": []string{"name"},
			},
			"page":  1,
			"limit": 30,
		},
	}
}

/**
* From
* @param model *Model
* @return *Ql
**/
func From(name interface{}) *Ql {
	var model *Model
	switch v := name.(type) {
	case *Model:
		model = v
	default:
		str := fmt.Sprintf("%v", v)
		model = GetModel(str)
	}

	tpSelect := Select
	if model.SourceField != nil {
		tpSelect = Source
	}

	result := &Ql{
		Id:         utility.UUID(),
		Db:         model.Db,
		TypeSelect: tpSelect,
		Froms:      newForms(),
		Joins:      make([]*QlJoin, 0),
		Selects:    make([]*Field, 0),
		Hiddens:    make([]string, 0),
		Details:    make([]*Field, 0),
		Groups:     make([]*Field, 0),
		Orders:     &QlOrder{Asc: make([]*Field, 0), Desc: make([]*Field, 0)},
		Offset:     0,
		Limit:      0,
		Sheet:      0,
		Help:       helpQl(),
		wg:         &sync.WaitGroup{},
	}
	result.QlWhere = newQlWhere(result.validator)
	result.IsDebug = model.IsDebug
	result.Havings = NewQlHaving(result)
	result.Froms.add(model)

	return result
}

/**
* From
* @param name string
* @return *Ql
**/
func (s *Ql) From(name string) *Ql {
	model := s.Db.GetModel(name)
	if model == nil {
		return s
	}

	main := s.Froms.add(model)
	for _, from := range s.Froms.Froms {
		if from.As != main.As {
			for _, detail := range from.RelationsTo {
				if detail.With.Id == main.Id {
					j := s.Join(main.Model)
					for fk, pk := range detail.Fk {
						j.On(fk).Eq(from.As + "." + pk)
					}
					return s
				}
			}
		}
	}

	return s
}

/**
* getForms
* @return []string
**/
func (s *Ql) getForms() []string {
	var result []string
	for _, from := range s.Froms.Froms {
		result = append(result, fmt.Sprintf(`%s.%s, %s`, from.Schema, from.Name, from.As))
	}

	return result
}
