package main

import (
	"github.com/cgalvisleon/et/cache"
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	_ "github.com/cgalvisleon/jdb/drivers/postgres"
	jdb "github.com/cgalvisleon/jdb/v1"
)

func main() {
	err := cache.Load()
	if err != nil {
		console.Panic(err)
	}

	err = event.Load()
	if err != nil {
		console.Panic(err)
	}

	// db, err := jdb.Load()
	// if err != nil {
	// 	console.Panic(err)
	// }

	model, err := jdb.DefineModel(et.Json{
		"database": "josephine",
		"schema":   "projects",
		"name":     "users",
		"version":  1,
		"columns": et.Json{
			"id": et.Json{
				"type":    "key",
				"default": "-1",
			},
			"name": et.Json{
				"type": "text",
			},
			"email": et.Json{
				"type": "text",
			},
		},
		"atribs": et.Json{
			"apellido": "", //name, defaultValue
			"rol":      "",
		},
		"primary_keys": []string{
			"id",
		},
		"details": et.Json{
			"roles": et.Json{
				"schema": "projects",
				"name":   "roles",
				"references": et.Json{
					"columns": et.Json{
						"user_id": "id",
					},
					"on_delete": "",
					"on_update": "",
				},
			},
		},
		"required": []string{"id"},
		"debug":    true,
	})
	if err != nil {
		console.Panic(err)
	}

	err = model.Init()
	if err != nil {
		console.Panic(err)
	}

	query, err := jdb.Query(et.Json{
		"database": "josephine",
		"from":     "users AS a",
		"select": []string{
			"a.id",
			"a.name",
			"a.email",
		},
		"joins": []et.Json{
			{
				"from": "roles AS b",
				"on":   "a.id = b.user_id",
			},
		},
		"where": et.Json{
			"a.id": et.Json{
				"eq": 1,
			},
		},
		"and": et.Json{
			"a.name": et.Json{
				"eq": "John",
			},
		},
		"or": et.Json{
			"a.name": et.Json{
				"eq": "Jane",
			},
		},
		"group_by": []string{""},
		"having": et.Json{
			"a.name": et.Json{
				"eq": "",
			},
		},
		"order_by": et.Json{
			"asc":  []string{""},
			"desc": []string{""},
		},
		"limit": et.Json{
			"page": 1,
			"rows": 10,
		},
	})
	if err != nil {
		console.Panic(err)
	}

	console.Debug("query:", query.ToJson().ToString())

	// result, err := query.ItExists()
	// if err != nil {
	// 	console.Panic(err)
	// }

	// console.Debug("exists:", result)
}
