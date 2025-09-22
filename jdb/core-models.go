package jdb

import (
	"github.com/cgalvisleon/et/et"
)

var models *Model

/**
* defineModel
* @param db *Database
* @return error
**/
func defineModel(db *Database) error {
	if models != nil {
		return nil
	}

	var err error
	models, err = db.Define(et.Json{
		"schema":  "core",
		"name":    "models",
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
				"name": "name",
				"type": "text",
			},
			{
				"name": "version",
				"type": "int",
			},
			{
				"name": "definition",
				"type": "bytes",
			},
		},
		"debug": true,
	})
	if err != nil {
		return err
	}

	err = models.Init()
	if err != nil {
		return err
	}

	return nil
}

/**
* setModel
* @param id string, data et.Json
* @return error
**/
func setModel(id string, data et.Json, debug bool) error {

	return nil
}

/**
* loadModel
* @param id string
* @return et.Json
**/
func loadModel(id string) et.Json {
	return et.Json{}
}
