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
		"fields": et.Json{
			"IndexField":     "index",
			"SourceField":    "_data",
			"ProjectField":   "project_id",
			"CreatedAtField": "created_at",
			"UpdatedAtField": "update_at",
			"StateField":     "_state",
			"KeyField":       "_id",
			"SystemKeyField": "_idt",
			"ClassField":     "_class",
			"CreatedToField": "created_to",
			"UpdatedToField": "updated_to",
			"FullTextField":  "_fulltext",
		},
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

	fields := params.Json("fields")
	for key, value := range fields {
		switch key {
		case "IndexField":
			IndexField = value.(ColumnField)
		case "SourceField":
			SourceField = value.(ColumnField)
		case "ProjectField":
			ProjectField = value.(ColumnField)
		case "CreatedAtField":
			CreatedAtField = value.(ColumnField)
		case "UpdatedAtField":
			UpdatedAtField = value.(ColumnField)
		case "StateField":
			StateField = value.(ColumnField)
		case "KeyField":
			KeyField = value.(ColumnField)
		case "SystemKeyField":
			SystemKeyField = value.(ColumnField)
		case "ClassField":
			ClassField = value.(ColumnField)
		case "CreatedToField":
			CreatedToField = value.(ColumnField)
		case "UpdatedToField":
			UpdatedToField = value.(ColumnField)
		}
	}

	dbs[name] = result

	return result, nil
}

func init() {
	drivers = map[string]func() Driver{}
	dbs = map[string]*DB{}
	schemas = map[string]*Schema{}
	models = map[string]*Model{}
}
