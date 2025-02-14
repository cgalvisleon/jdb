package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
)

type JDB struct {
	Drivers map[string]func() Driver
	DBs     map[string]*DB
	Schemas map[string]*Schema
	Models  map[string]*Model
	Flows   map[string]*Flow
	Version string
}

var Jdb *JDB

func init() {
	Jdb = &JDB{
		Drivers: map[string]func() Driver{},
		DBs:     map[string]*DB{},
		Schemas: map[string]*Schema{},
		Models:  map[string]*Model{},
		Flows:   map[string]*Flow{},
		Version: "0.0.1",
	}
}

/**
* Describe
* @return et.Json
**/
func (s *JDB) Describe() et.Json {
	drivers := []string{}
	for key := range s.Drivers {
		drivers = append(drivers, key)
	}
	dbs := []string{}
	for key := range s.DBs {
		dbs = append(dbs, key)
	}
	schemas := []et.Json{}
	for _, val := range s.Schemas {
		schemas = append(schemas, val.Describe())
	}
	models := []et.Json{}
	for _, val := range s.Models {
		models = append(models, val.Describe())
	}
	flows := []et.Json{}
	for _, val := range s.Flows {
		flows = append(flows, val.Describe())
	}

	result := et.Json{
		"drivers": drivers,
		"dbs":     dbs,
		"schemas": schemas,
		"models":  models,
		"flows":   flows,
		"version": s.Version,
	}

	return result
}

/**
* Load
* @return *DB, error
**/
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

/**
* ConnectTo
* @param params et.Json
* @return *DB, error
**/
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
		case "PrimaryKeyField":
			PrimaryKeyField = value.(ColumnField)
		case "SystemKeyField":
			SystemKeyField = value.(ColumnField)
		case "CreatedToField":
			CreatedToField = value.(ColumnField)
		case "UpdatedToField":
			UpdatedToField = value.(ColumnField)
		}
	}

	Jdb.DBs[name] = result

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
		return Jdb.Schemas[strs.Lowcase(name)]
	case 2:
		schema := Jdb.Schemas[list[1]]
		if schema != nil {
			return schema
		}
		if isCreate {
			db := Jdb.DBs[list[0]]
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
func GetModel(name string) *Model {
	list := strs.Split(name, ".")
	switch len(list) {
	case 1:
		for _, model := range Jdb.Models {
			if model.Name == strs.Lowcase(name) {
				return model
			}
		}
	case 2:
		table := strs.Format(`%s.%s`, list[0], list[1])
		result := Jdb.Models[strs.Lowcase(table)]
		if result != nil {
			return result
		}
	}

	return nil
}

/**
* GetField
* @param name string
* @return *Field
**/
func GetField(name string) *Field {
	list := strs.Split(name, ".")
	switch len(list) {
	case 2:
		model := GetModel(list[0])
		if model == nil {
			return nil
		}
		return model.GetField(list[1])
	case 3:
		table := strs.Format(`%s.%s`, list[0], list[1])
		model := GetModel(table)
		if model == nil {
			return nil
		}
		return model.GetField(list[2])
	default:
		return nil
	}
}

/**
* Describe
* @return et.Json
**/
func Describe() et.Json {
	return Jdb.Describe()
}
