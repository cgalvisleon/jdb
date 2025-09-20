package jdb

import (
	"github.com/cgalvisleon/et/console"
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
	models, err = db.DefineModel(et.Json{
		"schema": "congo",
		"name":   "models",
	})
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
	if debug {
		console.Debugf("%s:%s", id, data.ToString())
	}

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
