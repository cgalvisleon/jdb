package jdb

import "github.com/cgalvisleon/et/et"

var recycling *Model

/**
* defineRecycling
* @param db *Database
* @return error
**/
func defineRecycling(db *Database) error {
	if recycling != nil {
		return nil
	}

	var err error
	recycling, err = db.Define(et.Json{
		"schema":  "core",
		"name":    "recycling",
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
				"name": RECORDID,
				"type": "key",
			},
		},
		"record_field": RECORDID,
		"primary_keys": []string{RECORDID},
		"indices":      []string{},
		"debug":        true,
	})
	if err != nil {
		return err
	}

	recycling.isCore = true
	err = recycling.Init()
	if err != nil {
		return err
	}

	return nil
}
