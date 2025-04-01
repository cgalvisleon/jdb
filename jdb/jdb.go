package jdb

import (
	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/strs"
)

type JDB struct {
	Drivers map[string]func() Driver
	DBS     map[string]*DB
	Schemas map[string]*Schema
	Models  map[string]*Model
	Flows   map[string]*Flow
	Version string
}

var Jdb *JDB

func init() {
	Jdb = &JDB{
		Drivers: map[string]func() Driver{},
		DBS:     map[string]*DB{},
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
	for key := range s.DBS {
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
* ConnectTo
* @param params et.Json
* @return *DB, error
**/
func ConnectTo(params ConnectParams) (*DB, error) {
	driver := params.Driver
	if driver == "" {
		return nil, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	name := params.Name
	nodeId := params.NodeId
	result, err := NewDatabase(name, driver, nodeId)
	if err != nil {
		return nil, err
	}

	err = result.Conected(params.Params)
	if err != nil {
		return nil, err
	}

	if params.Fields != nil {
		for key, value := range params.Fields {
			switch key {
			case "IndexField":
				IndexField = ColumnField(value)
			case "SourceField":
				SourceField = ColumnField(value)
			case "ProjectField":
				ProjectField = ColumnField(value)
			case "CreatedAtField":
				CreatedAtField = ColumnField(value)
			case "UpdatedAtField":
				UpdatedAtField = ColumnField(value)
			case "StateField":
				StateField = ColumnField(value)
			case "PrimaryKeyField":
				PrimaryKeyField = ColumnField(value)
			case "SystemKeyField":
				SystemKeyField = ColumnField(value)
			case "CreatedToField":
				CreatedToField = ColumnField(value)
			case "UpdatedToField":
				UpdatedToField = ColumnField(value)
			case "FullTextField":
				FullTextField = ColumnField(value)
			}
		}
	}

	result.UseCore = params.UserCore
	if result.UseCore {
		err := result.CreateCore()
		if err != nil {
			return nil, err
		}
	}

	Jdb.DBS[name] = result

	return result, nil
}

/**
* Load
* @return *DB, error
**/
func Load() (*DB, error) {
	return ConnectTo(ConnectParams{
		Driver: envar.GetStr(PostgresDriver, "DB_DRIVER"),
		Name:   envar.GetStr("data", "DB_NAME"),
		NodeId: envar.GetInt64(0, "NODEID"),
		Params: et.Json{
			"database": envar.GetStr("data", "DB_NAME"),
			"host":     envar.GetStr("localhost", "DB_HOST"),
			"port":     envar.GetInt(5432, "DB_PORT"),
			"username": envar.GetStr("", "DB_USER"),
			"password": envar.GetStr("", "DB_PASSWORD"),
			"app":      envar.GetStr("jdb", "APP_NAME"),
		},
		Fields: map[string]string{
			"IndexField":     INDEX,
			"SourceField":    SOURCE,
			"ProjectField":   PROJECT,
			"CreatedAtField": CREATED_AT,
			"UpdatedAtField": UPDATED_AT,
			"StateField":     STATUS,
			"KeyField":       PRIMARYKEY,
			"SystemKeyField": SYSID,
			"CreatedToField": CREATED_TO,
			"UpdatedToField": UPDATED_TO,
			"FullTextField":  FULLTEXT,
		},
		UserCore: true,
	})
}

/**
* GetDB
* @param name string
* @return *DB
**/
func GetDB(name string) *DB {
	return Jdb.DBS[name]
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
		return Jdb.Schemas[name]
	case 2:
		schema := Jdb.Schemas[list[1]]
		if schema != nil {
			return schema
		}
		if isCreate {
			db := Jdb.DBS[list[0]]
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
* @param key string, isCreated bool
* @return *Model
**/
func GetModel(key string) *Model {
	list := strs.Split(key, ".")
	switch len(list) {
	case 1:
		for _, model := range Jdb.Models {
			if model.Name == key {
				return model
			}
		}
	case 2:
		result := Jdb.Models[key]
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
* @param name string
* @return et.Json
**/
func Describe(name string) (et.Json, error) {
	list := strs.Split(name, ":")
	if len(list) == 2 {
		prefix := list[0]
		switch prefix {
		case "db":
			db := GetDB(list[1])
			if db != nil {
				return db.Describe(), nil
			}

			return et.Json{}, mistake.Newf(MSG_DATABASE_NOT_FOUND, list[1])
		case "schema":
			sch := GetShema(list[1], false)
			if sch != nil {
				return sch.Describe(), nil
			}

			return et.Json{}, mistake.Newf(MSG_SCHEMA_NOT_FOUND, list[1])
		case "model":
			mod := GetModel(list[1])
			if mod != nil {
				return mod.Describe(), nil
			}

			return et.Json{}, mistake.Newf(MSG_MODEL_NOT_FOUND, list[1])
		}
	}

	mod := GetModel(name)
	if mod == nil {
		sch := GetShema(name, false)
		if sch == nil {
			result := et.Json{}
			for _, db := range JDBS {
				result.Set(db.Name, db.Describe())
			}

			return result, nil
		}

		return sch.Describe(), nil
	}

	return mod.Describe(), nil
}

/**
* Query
* @param params et.Json
* @return interface{}, error
**/
func Query(params et.Json) (interface{}, error) {
	from := params.Str("from")
	model := GetModel(from)
	if model == nil {
		return nil, mistake.Newf(MSG_MODEL_NOT_FOUND, from)
	}

	return From(model).
		Query(params)
}

/**
* Query
* @param params et.Json
* @return interface{}, error
**/
func Commands(params et.Json) (interface{}, error) {
	insert := params.Str("insert")
	update := params.Str("update")
	delete := params.Str("delete")
	data := params.Json("data")
	where := params.Json("where")
	var conn *Command
	if insert != "" {
		model := GetModel(insert)
		if model == nil {
			return nil, mistake.Newf(MSG_MODEL_NOT_FOUND, insert)
		}

		conn = model.Insert(data)
	} else if update != "" {
		model := GetModel(update)
		if model == nil {
			return nil, mistake.Newf(MSG_MODEL_NOT_FOUND, update)
		}

		conn = model.Update(data).
			setWhere(where)
	} else if delete != "" {
		model := GetModel(delete)
		if model == nil {
			return nil, mistake.Newf(MSG_MODEL_NOT_FOUND, delete)
		}

		conn = model.Delete().
			setWhere(where)
	} else {
		return nil, mistake.New(MSG_COMMAND_NOT_FOUND)
	}

	return conn.Exec()
}
