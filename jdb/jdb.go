package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

var (
	drivers   map[string]func() Driver
	dbs       map[string]*DB
	schemas   map[string]*Schema
	models    map[string]*Model
	functions map[string]func(et.Json) (et.Json, error)
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
		"nodeId":   result.Node,
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
		case "CreatedToField":
			CreatedToField = value.(ColumnField)
		case "UpdatedToField":
			UpdatedToField = value.(ColumnField)
		}
	}

	dbs[name] = result

	return result, nil
}

/**
* GetShema
* @param name string
* @param isCreate bool
* @return *Schema
**/
func GetShema(name string, isCreate bool) *Schema {
	list := strs.Split(name, ".")
	switch len(list) {
	case 1:
		return schemas[strs.Lowcase(name)]
	case 2:
		schema := schemas[list[1]]
		if schema != nil {
			return schema
		}
		if isCreate {
			db := dbs[list[0]]
			if db == nil {
				return nil
			}

			result, err := NewSchema(db, list[1])
			if err != nil {
				return nil
			}

			return result
		}
	}

	return nil
}

/**
* GetModel
* @param name string
* @param isCreated bool
* @return *Model
**/
func GetModel(name string, isCreated bool) *Model {
	list := strs.Split(name, ".")
	switch len(list) {
	case 1:
		for _, model := range models {
			if model.Name == strs.Lowcase(name) {
				return model
			}
		}
	case 2:
		table := strs.Format(`%s.%s`, list[0], list[1])
		result := models[strs.Lowcase(table)]
		if result != nil {
			return result
		}
		schema := schemas[list[0]]
		if schema == nil {
			return nil
		}
		if isCreated {
			return NewModel(schema, table, 1)
		}
	}

	return nil
}

/**
* GetField
* @param name string, isCreated bool
* @return *Field
**/
func GetField(name string, isCreated bool) *Field {
	list := strs.Split(name, ".")
	switch len(list) {
	case 2:
		model := GetModel(list[0], isCreated)
		if model == nil {
			return nil
		}
		return model.GetField(list[1], isCreated)
	case 3:
		table := strs.Format(`%s.%s`, list[0], list[1])
		model := GetModel(table, isCreated)
		if model == nil {
			return nil
		}
		return model.GetField(list[2], isCreated)
	default:
		return nil
	}
}

func init() {
	drivers = map[string]func() Driver{}
	dbs = map[string]*DB{}
	schemas = map[string]*Schema{}
	models = map[string]*Model{}
}
