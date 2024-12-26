package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
)

var (
	drivers map[string]func() Driver
	dbs     map[string]*DB
	schemas map[string]*Schema
	models  map[string]*Model
)

func Load() (*DB, error) {
	name := envar.GetStr("", "DB_NAME")
	result, err := NewDatabase(name, Postgres)
	if err != nil {
		return nil, err
	}

	result.UseCore = true
	err = result.Conected(et.Json{
		"driver":   Postgres,
		"host":     envar.GetStr("localhost", "DB_HOST"),
		"port":     envar.GetInt(5432, "DB_PORT"),
		"database": name,
		"username": envar.GetStr("", "DB_USER"),
		"password": envar.GetStr("", "DB_PASSWORD"),
		"app":      envar.GetStr("jdb", "DB_APP_NAME"),
		"core":     result.UseCore,
	})
	if err != nil {
		return nil, err
	}

	result.CreateCore()

	return result, nil
}

func ConnectTo(params et.Json) (*DB, error) {
	driver := params.Str("driver")
	if driver == "" {
		return nil, console.Alertm("Driver not defined")
	}

	name := params.ValStr("db", "name")
	result, err := NewDatabase(name, driver)
	if err != nil {
		return nil, err
	}

	err = result.Conected(params)
	if err != nil {
		return nil, err
	}

	core := params.Bool("core")
	if core {
		result.CreateCore()
	}

	return result, nil
}

func init() {
	drivers = map[string]func() Driver{}
	dbs = map[string]*DB{}
	schemas = map[string]*Schema{}
	models = map[string]*Model{}
}
