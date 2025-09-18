package main

import (
	"github.com/cgalvisleon/et/cache"
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/event"
	_ "github.com/cgalvisleon/jdb/drivers/postgres"
	"github.com/cgalvisleon/jdb/jdb"
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
		// "columns": et.Json{
		// 	"id": et.Json{
		// 		"type":    "key",
		// 		"default": "-1",
		// 	},
		// 	"name": et.Json{
		// 		"type": "text",
		// 	},
		// 	"email": et.Json{
		// 		"type": "text",
		// 	},
		// },
		// "primary_keys": []string{
		// 	"id",
		// },
		// "foreign_keys": et.Json{
		// 	"id": "id",
		// },
		// "indices": []string{
		// 	"id",
		// 	"name",
		// 	"email",
		// },
		// "uniques":  []string{"id"},
		// "required": []string{"id"},
		// "details": et.Json{
		// 	"profiles": et.Json{},
		// },
	})
	if err != nil {
		console.Panic(err)
	}

	ddl, err := model.Load()
	if err != nil {
		console.Panic(err)
	}

	console.Debug("ddl:", ddl)

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
