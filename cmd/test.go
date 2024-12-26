package main

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	_ "github.com/cgalvisleon/jdb/drivers/postgres"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func test() {
	db, err := jdb.Load()
	if err != nil {
		panic(err)
	}

	model, _ := jdb.NewSchema(db, "test")
	users := jdb.NewModel(model, "users", 1)
	users.DefineColumn("name", jdb.TypeDataText)
	users.DefineColumn("age", jdb.TypeDataInt)
	users.DefineColumn("email", jdb.TypeDataText)

	items, err := jdb.From(users).
		Where("name").Eq("Carlos").
		And("age").More(18).
		Or("email").Like("*gmail.com").
		Select("name", "email").
		Debug().
		All()
	if err != nil {
		panic(err)
	}

	console.Debug(items.ToString())

	item, err := users.Insert(et.Json{}).
		Debug().
		One()
	if err != nil {
		panic(err)
	}

	console.Debug(item.ToString())
}
