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
		"schema": "core",
		"name":   "models",
		"columns": et.Json{
			"created_at": et.Json{
				"type": "datetime",
			},
			"updated_at": et.Json{
				"type": "datetime",
			},
			"kind": et.Json{
				"type": "key",
			},
			"name": et.Json{
				"type": "text",
			},
			"version": et.Json{
				"type": "int",
			},
			"definition": et.Json{
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
