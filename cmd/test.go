package main

import (
	_ "github.com/cgalvisl/jdb/drivers/postgres"
	jdb "github.com/cgalvisl/jdb/pkg"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
)

func test() {
	db, err := jdb.NewDatabase("test", "postgres")
	if err != nil {
		panic(err)
	}

	model := jdb.NewSchema(db, "test")
	users := jdb.NewModel(model, "users")
	users.DefineColumn("name", jdb.TypeDataText, "")
	users.DefineColumn("age", jdb.TypeDataInt, 0)
	users.DefineColumn("email", jdb.TypeDataText, "")

	items, err := jdb.From(users).
		Where("name").Eq("Carlos").
		And("age").More(18).
		Or("email").Like("*gmail.com").
		Data("name", "email").
		Debug().
		All()
	if err != nil {
		panic(err)
	}

	logs.Debug(items.ToString())

	item, err := users.Insert(et.Json{}).
		Debug().
		One()
	if err != nil {
		panic(err)
	}

	logs.Debug(item.ToString())
}
