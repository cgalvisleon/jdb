package jdb

import (
	"encoding/json"
	"slices"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/strs"
)

type JDB struct {
	Drivers map[string]func() Driver `json:"-"`
	DBS     []*DB                    `json:"-"`
	Version string                   `json:"version"`
}

var conn *JDB

func init() {
	conn = &JDB{
		Drivers: map[string]func() Driver{},
		DBS:     make([]*DB, 0),
		Version: "0.0.1",
	}
}

/**
* Describe
* @return et.Json
**/
func (s *JDB) Describe() et.Json {
	definition, err := json.Marshal(s)
	if err != nil {
		return et.Json{}
	}

	result := et.Json{}
	err = json.Unmarshal(definition, &result)
	if err != nil {
		return et.Json{}
	}

	drivers := []string{}
	for key := range s.Drivers {
		drivers = append(drivers, key)
	}

	dbs := []et.Json{}
	for _, db := range s.DBS {
		dbs = append(dbs, db.Describe())
	}

	result["drivers"] = drivers
	result["dbs"] = dbs

	return result
}

/**
* ConnectTo
* @param params et.Json
* @return *DB, error
**/
func ConnectTo(params ConnectParams) (*DB, error) {
	err := params.Validate()
	if err != nil {
		return nil, err
	}

	result, err := NewDatabase(params.Name, params.Driver)
	if err != nil {
		return nil, err
	}

	err = result.Conected(params.Params)
	if err != nil {
		return nil, err
	}

	result.UseCore = params.UserCore
	if result.UseCore {
		err := result.createCore()
		if err != nil {
			return nil, err
		}

		item, err := result.getModel("db", result.Name)
		if err != nil {
			return nil, err
		}

		if item.Ok {
			id := item.Str(SYSID)
			if id != "" {
				result.Id = id
			}
		}
	}

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
		Params: et.Json{
			"database": envar.GetStr("data", "DB_NAME"),
			"host":     envar.GetStr("localhost", "DB_HOST"),
			"port":     envar.GetInt(5432, "DB_PORT"),
			"username": envar.GetStr("", "DB_USER"),
			"password": envar.GetStr("", "DB_PASSWORD"),
			"app":      envar.GetStr("jdb", "APP_NAME"),
		},
		UserCore: true,
	})
}

/**
* Jdb
* @return *JDB
**/
func Jdb() *JDB {
	return conn
}

/**
* GetDB
* @param name string
* @return *DB
**/
func GetDB(name string) *DB {
	name = Name(name)
	idx := slices.IndexFunc(conn.DBS, func(db *DB) bool { return db.Name == name })
	if idx != -1 {
		return conn.DBS[idx]
	}

	return nil
}

/**
* GetShema
* @param name string
* @param isCreate bool
* @return *Schema
**/
func GetShema(name string, isCreate bool) *Schema {
	if len(conn.DBS) == 0 {
		return nil
	}

	name = Name(name)
	var db *DB
	var result *Schema
	list := strs.Split(name, ".")
	switch len(list) {
	case 1:
		db = conn.DBS[0]
		if db == nil {
			return nil
		}

		result = db.GetSchema(name)
		if result != nil {
			return result
		}
	case 2:
		db = GetDB(list[0])
		if db == nil {
			return nil
		}

		name := list[1]
		result = db.GetSchema(name)
		if result != nil {
			return result
		}
	default:
		return result
	}

	if isCreate {
		result = NewSchema(db, name)
	}

	return result
}

/**
* GetModel
* @param name string, isCreate bool
* @return *Model
**/
func GetModel(name string, isCreate bool) *Model {
	if len(conn.DBS) == 0 {
		return nil
	}

	var db *DB
	var schema *Schema
	var result *Model
	list := strs.Split(name, ".")
	switch len(list) {
	case 1: // model
		db = conn.DBS[0]
		if db == nil {
			return nil
		}

		result = db.GetModel(name)
		if result != nil {
			return result
		}
	case 2: // schema, model
		db = conn.DBS[0]
		if db == nil {
			return nil
		}

		schema = db.GetSchema(list[0])
		if schema == nil {
			return nil
		}

		name := list[1]
		result = schema.GetModel(name)
		if result != nil {
			return result
		}
	case 3: // db, schema, model
		db = GetDB(list[0])
		if db == nil {
			return nil
		}

		schema = db.GetSchema(list[1])
		if schema == nil {
			return nil
		}

		name := list[2]
		result = schema.GetModel(name)
		if result != nil {
			return result
		}
	default:
		return result
	}

	if !isCreate {
		return result
	}

	if schema == nil {
		return nil
	}

	return NewModel(schema, name, 0)
}

/**
* Describe
* @param name string
* @return et.Json
**/
func Describe(name string) (et.Json, error) {
	list := strs.Split(name, ":")
	switch len(list) {
	case 1: // model
		model := GetModel(name, false)
		if model != nil {
			return model.Describe(), nil
		}

		schema := GetShema(name, false)
		if schema != nil {
			return schema.Describe(), nil
		}

		db := GetDB(name)
		if db != nil {
			return db.Describe(), nil
		}

		return conn.Describe(), nil
	case 2: // schema, model
		schema := GetShema(list[0], false)
		if schema == nil {
			return et.Json{}, mistake.Newf(MSG_SCHEMA_NOT_FOUND, list[0])
		}

		model := schema.GetModel(list[1])
		if model != nil {
			return model.Describe(), nil
		}

		return schema.Describe(), nil
	case 3: // db, schema, model
		db := GetDB(list[0])
		if db == nil {
			return et.Json{}, mistake.Newf(MSG_DATABASE_NOT_FOUND, list[0])
		}

		schema := db.GetSchema(list[1])
		if schema == nil {
			return et.Json{}, mistake.Newf(MSG_SCHEMA_NOT_FOUND, list[1])
		}

		model := schema.GetModel(list[2])
		if model != nil {
			return model.Describe(), nil
		}

		return schema.Describe(), nil
	case 4: // db, schema, model, field
		db := GetDB(list[0])
		if db == nil {
			return et.Json{}, mistake.Newf(MSG_DATABASE_NOT_FOUND, list[0])
		}

		schema := db.GetSchema(list[1])
		if schema == nil {
			return et.Json{}, mistake.Newf(MSG_SCHEMA_NOT_FOUND, list[1])
		}

		model := schema.GetModel(list[2])
		if model == nil {
			return et.Json{}, mistake.Newf(MSG_MODEL_NOT_FOUND, list[2])
		}

		field := model.getField(list[3], false)
		if field != nil {
			return field.Describe(), nil
		}

		return model.Describe(), nil
	}

	return et.Json{}, mistake.Newf(MSG_INVALID_NAME, name)
}

/**
* Query
* @param params et.Json
* @return interface{}, error
**/
func Query(params et.Json) (interface{}, error) {
	from := params.Str("from")
	model := GetModel(from, false)
	if model == nil {
		return nil, mistake.Newf(MSG_MODEL_NOT_FOUND, from)
	}

	return From(model).
		Query(params)
}

/**
* Commands
* @param params et.Json
* @return interface{}, error
**/
func Commands(params et.Json) (interface{}, error) {
	insert := params.Str("insert")
	update := params.Str("update")
	delete := params.Str("delete")
	data := params.Json("data")
	where := params.Json("where")
	var comm *Command
	if insert != "" {
		model := GetModel(insert, false)
		if model == nil {
			return nil, mistake.Newf(MSG_MODEL_NOT_FOUND, insert)
		}

		comm = model.Insert(data)
	} else if update != "" {
		model := GetModel(update, false)
		if model == nil {
			return nil, mistake.Newf(MSG_MODEL_NOT_FOUND, update)
		}

		comm = model.Update(data).
			setWhere(where)
	} else if delete != "" {
		model := GetModel(delete, false)
		if model == nil {
			return nil, mistake.Newf(MSG_MODEL_NOT_FOUND, delete)
		}

		comm = model.Delete().
			setWhere(where)
	} else {
		return nil, mistake.New(MSG_COMMAND_NOT_FOUND)
	}

	return comm.Exec()
}
