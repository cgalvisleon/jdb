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

	// user := jdb.NewModel(db.Schema, "user", 1)
	// user.DefineProjectModel()
	// result, err := user.
	// 	Insert(et.Json{}).
	// 	BeforeInsert(func(tx *jdb.Tx, data et.Json) error {
	// 		data["created_at"] = time.Now()
	// 		data["updated_at"] = time.Now()
	// 		return nil
	// 	}).
	// 	BeforeUpdate(func(tx *jdb.Tx, data et.Json) error {
	// 		data["updated_at"] = time.Now()
	// 		return nil
	// 	}).
	// 	BeforeDelete(func(tx *jdb.Tx, data et.Json) error {
	// 		data["deleted_at"] = time.Now()
	// 		return nil
	// 	}).
	// 	BeforeInsertOrUpdate(func(tx *jdb.Tx, data et.Json) error {
	// 		data["updated_at"] = time.Now()
	// 		return nil
	// 	}).
	// 	AfterInsert(func(tx *jdb.Tx, data et.Json) error {
	// 		console.Debug("after insert", data)
	// 		return nil
	// 	}).
	// 	AfterUpdate(func(tx *jdb.Tx, data et.Json) error {
	// 		console.Debug("after update", data)
	// 		return nil
	// 	}).
	// 	AfterDelete(func(tx *jdb.Tx, data et.Json) error {
	// 		console.Debug("after delete", data)
	// 		return nil
	// 	}).
	// 	AfterInsertOrUpdate(func(tx *jdb.Tx, data et.Json) error {
	// 		console.Debug("after insert or update", data)
	// 		return nil
	// 	}).
	// 	Exec()
	// if err != nil {
	// 	panic(err)
	// }

	// result := db.Describe()
	// console.Debug(result.ToString())
	console.Debug("db:", db.Name)
}
