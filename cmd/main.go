package main

import "github.com/cgalvisl/jdb/jdb"

func main() {

	db := jdb.NewDatabase("test", "")
	model := jdb.NewSchema(db, "test", "")
	users := jdb.NewModel(model, "users", "")
	users.DefineColumn("name", "", jdb.TypeDataText, "")
	users.DefineColumn("age", "", jdb.TypeDataInt, 0)
	users.DefineColumn("email", "", jdb.TypeDataText, "")

	jdb.From(users).
		Where("name").
		Eq("Carlos").
		And("age").
		More(18).
		Or("email").
		Like("*gmail.com").
		Select().
		All()
}
