package main

import (
	"github.com/cgalvisleon/et/console"
	_ "github.com/cgalvisleon/jdb/drivers/postgres"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func main() {
	result, err := jdb.Describe("")
	if err != nil {
		panic(err)
	}

	console.Debug(result.ToString())
}
