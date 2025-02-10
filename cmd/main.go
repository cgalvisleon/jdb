package main

import (
	"github.com/cgalvisleon/et/console"
	_ "github.com/cgalvisleon/jdb/drivers/postgres"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func main() {
	// _, err := jdb.Load()
	// if err != nil {
	// 	panic(err)
	// }

	console.Debug(jdb.Describe().ToString())
}
