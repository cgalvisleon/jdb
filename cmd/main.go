package main

import (
	"github.com/cgalvisleon/et/cache"
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/event"
	"github.com/cgalvisleon/et/reg"
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

	db, err := jdb.LoadTo("josephine")
	if err != nil {
		console.Panic(err)
	}

	id := reg.GenULIDI("users")
	console.Debug("id:", id)

	console.Debug("db:", db.ToJson().ToString())

	// model := db.Models["users"]
	// if model == nil {
	// 	console.Panic(fmt.Errorf("model not found"))
	// }

	// result, err := model.Where(jdb.Eq("_id", "USERA00000001")).
	// 	Join("project_users", "B", jdb.Eq("A._id", "B.user_id")).
	// 	And(jdb.Eq("A.name", "John")).
	// 	Group("A._id", "A.caption").
	// 	Having(jdb.Eq("A._id", "USERA00000001")).
	// 	Order(true, "A._id", "A.caption").
	// 	Limit(1, 10)
	// if err != nil {
	// 	console.Panic(err)
	// }

	// console.Debug("result:", result.ToString())

	// users, err := db.Query(`SELECT json_build_object(
	// 	'_id', _id,
	// 	'username', username,
	// 	'name', caption
	// ) AS result FROM js_core.users`)
	// if err != nil {
	// 	console.Panic(err)
	// }

	// console.Debug("users:", users.ToString())

	// model, err := db.Define(et.Json{
	// 	"schema":  "projects",
	// 	"name":    "users",
	// 	"version": 1,
	// 	"columns": []et.Json{
	// 		{
	// 			"name":    "id",
	// 			"type":    "key",
	// 			"default": "-1",
	// 		},
	// 		{
	// 			"name": "name",
	// 			"type": "text",
	// 		},
	// 		{
	// 			"name": "email",
	// 			"type": "text",
	// 		},
	// 	},
	// 	"atribs": et.Json{
	// 		"apellido": "", //name, defaultValue
	// 		"rol":      "",
	// 	},
	// 	"primary_keys": []string{
	// 		"id",
	// 	},
	// 	"details": []et.Json{
	// 		{
	// 			"schema": "projects",
	// 			"name":   "roles",
	// 			"references": et.Json{
	// 				"columns": []et.Json{
	// 					{
	// 						"user_id": "id",
	// 					},
	// 				},
	// 				"on_delete": "",
	// 				"on_update": "",
	// 			},
	// 		},
	// 	},
	// 	"required": []string{"id"},
	// 	"debug":    true,
	// })
	// if err != nil {
	// 	console.Panic(err)
	// }

	// err = model.Init()
	// if err != nil {
	// 	console.Panic(err)
	// }

	// query, err := db.Select(et.Json{
	// 	"source_field": "_data",
	// 	"from": et.Json{
	// 		"js_core.users": "A",
	// 	},
	// 	"select": et.Json{
	// 		"A._id":     "_id",
	// 		"A.caption": "caption",
	// 	},
	// 	"atribs": et.Json{
	// 		"_data#>>'{created_by}'": "created_by",
	// 	},
	// 	"joins": []et.Json{
	// 		{
	// 			"from": et.Json{
	// 				"js_core.project_users": "B",
	// 			},
	// 			"on": []et.Json{
	// 				{
	// 					"A._id": et.Json{
	// 						"eq": "B.user_id",
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	"where": []et.Json{
	// 		{
	// 			"A._id": et.Json{
	// 				"eq": "'USERA00000001'",
	// 			},
	// 		},
	// 		{
	// 			"and": []et.Json{
	// 				{
	// 					"A.name": et.Json{
	// 						"eq": "'John'",
	// 					},
	// 				},
	// 			},
	// 		},
	// 		{
	// 			"or": []et.Json{
	// 				{
	// 					"A.name": et.Json{
	// 						"eq": "'Jane'",
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	"group_by": []string{"A.uno", "A.dos"},
	// 	"having": []et.Json{
	// 		{
	// 			"A.name": et.Json{
	// 				"eq": "'col'",
	// 			},
	// 		},
	// 	},
	// 	"order_by": et.Json{
	// 		"asc":  []string{"A.name", "A.apellidos"},
	// 		"desc": []string{"B.apellidos", "B.name"},
	// 	},
	// 	"limit": et.Json{
	// 		"page": 1,
	// 		"rows": 10,
	// 	},
	// })
	// if err != nil {
	// 	console.Panic(err)
	// }

	// console.Debug("query:", query.ToJson().ToString())

	// result, err := query.All()
	// if err != nil {
	// console.Panic(err)
	// }

	// console.Debug("exists:", result.ToString())
}
