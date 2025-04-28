package jdb

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/response"
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

type ConnectParams struct {
	Driver   string  `json:"driver"`
	Name     string  `json:"name"`
	Params   et.Json `json:"params"`
	UserCore bool    `json:"user_core"`
}

/**
* Validate
* @return error
**/
func (s *ConnectParams) validate() error {
	if conn == nil {
		return mistake.New(MSG_JDB_NOT_DEFINED)
	}

	if s.Driver == "" {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	if s.Name == "" {
		return mistake.New(MSG_DATABASE_NOT_DEFINED)
	}

	if _, ok := conn.Drivers[s.Driver]; !ok {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return nil
}

/**
* Serialize
* @return []byte, error
**/
func (s *JDB) Serialize() ([]byte, error) {
	result, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}

	return result, nil
}

/**
* Describe
* @return et.Json
**/
func (s *JDB) Describe() et.Json {
	definition, err := s.Serialize()
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
	err := params.validate()
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
* @param kind, name string
* @return et.Json
**/
func Describe(kind, name string) (et.Json, error) {
	help := et.Json{
		"message": MSG_KIND_NOT_DEFINED,
		"help":    "Exist four types of objects: db, schema, model and field. It is required at least two params, kind and name.",
		"params": et.Json{
			"kind": "db",
			"name": "name",
		},
	}
	if kind == "" {
		return help, mistake.New(MSG_KIND_NOT_DEFINED)
	}

	switch kind {
	case "db":
		result := GetDB(name)
		if result == nil {
			return et.Json{}, mistake.Newf(MSG_DATABASE_NOT_FOUND, name)
		}

		return result.Describe(), nil
	case "schema":
		result := GetShema(name, false)
		if result == nil {
			return et.Json{}, mistake.Newf(MSG_SCHEMA_NOT_FOUND, name)
		}

		return result.Describe(), nil
	case "model":
		result := GetModel(name, false)
		if result == nil {
			return et.Json{}, mistake.Newf(MSG_MODEL_NOT_FOUND, name)
		}

		return result.Describe(), nil
	case "field":
		list := strs.Split(name, ".")
		if len(list) != 2 {
			return et.Json{
				"message": MSG_INVALID_NAME,
				"help":    "It is required at least two parts in the name of the field, first part is the name of model and second is field name.",
				"example": "model.field",
			}, mistake.New(MSG_INVALID_NAME)
		}

		model := GetModel(list[0], false)
		if model == nil {
			return et.Json{}, mistake.Newf(MSG_MODEL_NOT_FOUND, list[0])
		}

		field := model.getField(list[1], false)
		if field == nil {
			return et.Json{}, mistake.Newf(MSG_FIELD_NOT_FOUND, list[1])
		}

		return field.Describe(), nil
	}

	return help, nil
}

/**
* Define
* @param params et.Json
* @return et.Json, error
**/
func Define(params et.Json) (et.Json, error) {
	result := et.Json{}
	help := et.Json{
		"message": MSG_INVALID_MODEL_PARAM,
		"help":    "It is required this params.",
		"params": et.Json{
			"name_model": et.Json{
				"schema":  "schema_name",
				"version": 1,
				"fields":  []et.Json{},
			},
		},
	}
	for name := range params {
		param := params.Json(name)
		if param.IsEmpty() {
			return help, mistake.Newf(MSG_INVALID_MODEL_PARAM, name)
		}

		schemaName := param.Str("schema")
		schema := GetShema(schemaName, true)
		if schema == nil {
			return et.Json{}, mistake.Newf(MSG_SCHEMA_NOT_FOUND, schemaName)
		}

		version := param.Int("version")
		model := NewModel(schema, name, version)
		if model == nil {
			return et.Json{}, mistake.Newf(MSG_MODEL_NOT_FOUND, name)
		}

		err := model.Save()
		if err != nil {
			return et.Json{}, err
		}

		result[name] = model.Describe()
	}

	return result, nil
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
			setWheres(where)
	} else if delete != "" {
		model := GetModel(delete, false)
		if model == nil {
			return nil, mistake.Newf(MSG_MODEL_NOT_FOUND, delete)
		}

		comm = model.Delete().
			setWheres(where)
	} else {
		return nil, mistake.New(MSG_COMMAND_NOT_FOUND)
	}

	return comm.Exec()
}

/**
* ModelDescribe
* @param w http.ResponseWriter
* @param r *http.Request
**/
func ModelDescribe(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	kind := body.Str("kind")
	name := body.Str("name")
	result, err := Describe(kind, name)
	if err != nil {
		response.JSON(w, r, http.StatusBadRequest, result)
		return
	}

	response.JSON(w, r, http.StatusOK, result)
}

/**
* modelDefine
* @param w http.ResponseWriter
* @param r *http.Request
**/
func ModelDefine(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	params := body.Json("params")
	result, err := Define(params)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, r, http.StatusOK, result)
}

/**
* modelQuery
* @param w http.ResponseWriter
* @param r *http.Request
**/
func ModelQuery(w http.ResponseWriter, r *http.Request) {
	params, _ := response.GetBody(r)
	result, err := Query(params)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ANY(w, r, http.StatusOK, result)
}

/**
* modelCommand
* @param w http.ResponseWriter
* @param r *http.Request
**/
func ModelCommand(w http.ResponseWriter, r *http.Request) {
	params, _ := response.GetBody(r)
	result, err := Commands(params)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ANY(w, r, http.StatusOK, result)
}
