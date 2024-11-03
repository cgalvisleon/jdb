package main

import "github.com/cgalvisl/jdb/jdb"

func main() {

	db := jdb.NewDatabase("test", "")
	model := jdb.NewSchema(db, "test", "")
	users := jdb.NewModel(model, "users", "")

	jdb.From(users).
		Where().
		And().
		Or().
		Select().
		All()
}
