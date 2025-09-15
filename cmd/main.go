package main

import (
	"github.com/cgalvisleon/et/cache"
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/event"
	_ "github.com/cgalvisleon/jdb/drivers/postgres"
	jdb "github.com/cgalvisleon/jdb/jdb"
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

	db, err := jdb.Load()
	if err != nil {
		console.Panic(err)
	}

	console.Debug("db:", db.Name)
}
