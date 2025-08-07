package main

import (
	"github.com/cgalvisleon/et/cache"
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/event"
	_ "github.com/cgalvisleon/jdb/drivers/sqlite"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func main() {
	err := cache.Load()
	if err != nil {
		panic(err)
	}

	err = event.Load()
	if err != nil {
		panic(err)
	}

	db, err := jdb.Load()
	if err != nil {
		panic(err)
	}

	console.Debug("db:", db.Name)
}
