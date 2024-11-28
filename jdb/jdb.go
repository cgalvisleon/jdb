package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
)

var (
	drivers map[string]func() Driver
	dbs     map[string]*DB
	schemas map[string]*Schema
	models  map[string]*Model
)

func Load() (*DB, error) {
	result, err := NewDatabase(Postgres)
	if err != nil {
		return nil, err
	}

	err = result.Conected(et.Json{})
	if err != nil {
		return nil, err
	}

	logs.Log(Postgres, "Database connected")

	return result, nil
}

func ConnectTo(params et.Json) (*DB, error) {
	driver := params.Str("driver")
	if driver == "" {
		return nil, logs.NewError("Driver not defined")
	}

	result, err := NewDatabase(driver)
	if err != nil {
		return nil, err
	}

	err = result.Conected(params)
	if err != nil {
		return nil, err
	}

	logs.Log(driver, "Database connected")

	return result, nil
}

func init() {
	drivers = map[string]func() Driver{}
	dbs = map[string]*DB{}
	schemas = map[string]*Schema{}
	models = map[string]*Model{}
}
