package jdb

import (
	"github.com/cgalvisleon/et/envar"
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

	err = result.Conected(et.Json{
		"driver":   Postgres,
		"host":     envar.GetStr("localhost", "DB_HOST"),
		"port":     envar.GetInt(5432, "DB_PORT"),
		"database": envar.GetStr("", "DB_NAME"),
		"username": envar.GetStr("", "DB_USER"),
		"password": envar.GetStr("", "DB_PASSWORD"),
		"app":      envar.GetStr("jdb", "DB_APP_NAME"),
	})
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
