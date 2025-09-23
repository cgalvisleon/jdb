package jdb

import "github.com/cgalvisleon/et/et"

var series *Model

/**
* defineSeries
* @param db *Database
* @return error
**/
func defineSeries(db *Database) error {
	if series != nil {
		return nil
	}

	var err error
	series, err = db.Define(et.Json{
		"schema":  "core",
		"name":    "series",
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
				"name": "kind",
				"type": "key",
			},
			{
				"name": "tag",
				"type": "key",
			},
			{
				"name": "value",
				"type": "int",
			},
			{
				"name": "format",
				"type": "key",
			},
			{
				"name": RECORDID,
				"type": "key",
			},
		},
		"record_field": RECORDID,
		"primary_keys": []string{"kind", "tag"},
		"indices":      []string{RECORDID},
		"debug":        true,
	})
	if err != nil {
		return err
	}

	series.IsCore = true
	err = series.Init()
	if err != nil {
		return err
	}

	return nil
}
