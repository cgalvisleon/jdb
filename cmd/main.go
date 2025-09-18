package main

import (
	"github.com/cgalvisleon/et/cache"
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	jdb "github.com/cgalvisleon/jdb/congo"
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

	model, err := jdb.Define(et.Json{
		"database": "catalogo",
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
		// "foreign_keys": et.Json{
		// 	"id": "id",
		// },
		// "indices": []string{
		// 	"id",
		// 	"name",
		// 	"email",
		// },
		// "uniques":  []string{"id"},
		// 	"profiles": et.Json{},
		// },
	})
	if err != nil {
		console.Panic(err)
	}

	err = model.Init()
	if err != nil {
		console.Panic(err)
	}

	// sql, err := js.Query(et.Json{
	// 	"users": et.Json{
	// 		"select": []string{"id", "name", "email as correo"},
	// 		"where": et.Json{
	// 			"id": et.Json{
	// 				"eq": 1,
	// 			},
	// 		},
	// 		"and": et.Json{
	// 			"name": et.Json{
	// 				"eq": "John",
	// 			},
	// 			"profiles": et.Json{
	// 				"verified": et.Json{
	// 					"eq": true,
	// 				},
	// 			},
	// 		},
	// 	},
	// })
	// if err != nil {
	// 	console.Panic(err)
	// }

	// console.Debug("sql:", sql)

	// console.Debug("db:", db.Name)
}
