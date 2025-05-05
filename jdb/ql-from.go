package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/strs"
)

type QlFrom struct {
	*Model
	As string
}

type QlFroms struct {
	Froms []*QlFrom
	index int
}

func From(model *Model) *Ql {
	result := &Ql{
		Id:         reg.GenId("ql"),
		Db:         model.Db,
		TypeSelect: Select,
		Froms:      &QlFroms{index: 65, Froms: make([]*QlFrom, 0)},
		Joins:      make([]*QlJoin, 0),
		QlWhere:    NewQlWhere(),
		Selects:    make([]*Field, 0),
		Details:    make([]*Field, 0),
		Groups:     make([]*Field, 0),
		Orders:     &QlOrder{Asc: make([]*Field, 0), Desc: make([]*Field, 0)},
		Offset:     0,
		Limit:      0,
		Sheet:      0,
		Help: et.Json{
			"from": model.Name,
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
	result.Havings = &QlHaving{Ql: result, QlWhere: NewQlWhere()}
	result.addFrom(model)

	return result
}

/**
* listForms
* @return []string
**/
func (s *Ql) listForms() []string {
	var result []string
	for _, from := range s.Froms.Froms {
		result = append(result, strs.Format(`%s.%s, %s`, from.Schema, from.Name, from.As))
	}

	return result
}
