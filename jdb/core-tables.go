package jdb

import "github.com/cgalvisleon/et/et"

var tables *Model

/**
* defineTables
* @param db *Database
* @return error
**/
func defineTables(db *Database) error {
	if tables != nil {
		return nil
	}

	var err error
	tables, err = db.Define(et.Json{
		"schema":  "core",
		"name":    "tables",
		"version": 1,
		"columns": []et.Json{
			{
				"name": "created_at",
				"type": "datetime",
			},
			{
				"name": "updated_at",
				"type": "datetime",
			},
			{
				"name": "schema",
				"type": "key",
			},
			{
				"name": "table",
				"type": "text",
			},
			{
				"name": "total",
				"type": "int",
			},
		},
		"primary_keys": []string{"schema", "table"},
		"indices":      []string{},
		"debug":        true,
	})
	if err != nil {
		return err
	}

	tables.isCore = true
	err = tables.Init()
	if err != nil {
		return err
	}

	return nil
}
