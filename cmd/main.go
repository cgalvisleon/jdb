package main

import (
	"github.com/cgalvisleon/et/cache"
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/event"
	_ "github.com/cgalvisleon/jdb/drivers/postgres"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func main() {
	_, err := cache.Load()
	if err != nil {
		panic(err)
	}

	_, err = event.Load()
	if err != nil {
		panic(err)
	}

	db, err := jdb.Load()
	if err != nil {
		panic(err)
	}

	result := db.Describe()
	console.Debug(result.ToString())
}
