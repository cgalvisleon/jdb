package jdb

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/cgalvisleon/et/config"
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/response"
	"github.com/cgalvisleon/et/strs"
)

type JDB struct {
	Drivers   map[string]func() Driver `json:"-"`
	Params    map[string]ConnectParams `json:"-"`
	DBS       map[string]*DB           `json:"-"`
	DefaultDB *DB                      `json:"-"`
	Version   string                   `json:"version"`
}

var (
	os   = ""
	conn *JDB
)

func init() {
	os = runtime.GOOS
	conn = &JDB{
		Drivers: map[string]func() Driver{},
		Params:  map[string]ConnectParams{},
		DBS:     make(map[string]*DB),
		Version: "0.0.1",
	}
}

type ConnectParams struct {
	Driver   string   `json:"driver"`
	Name     string   `json:"name"`
	UserCore bool     `json:"user_core"`
	NodeId   int      `json:"node_id"`
	Debug    bool     `json:"debug"`
	Validate []string `json:"validate"`
	Params   et.Json  `json:"params"`
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

	result.IsDebug = params.Debug
	result.UseCore = params.UserCore
	result.NodeId = params.NodeId
	err = result.Conected(params.Params)
	if err != nil {
		return nil, err
	}

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
	driver := config.String("DB_DRIVER", SqliteDriver)
	if driver == "" {
		return nil, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	params := conn.Params[driver]
	err := config.Validate(params.Validate)
	if err != nil {
		return nil, err
	}

	return ConnectTo(params)
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
	if _, ok := conn.DBS[name]; ok {
		return conn.DBS[name]
	}

	if conn.DefaultDB != nil {
		return conn.DefaultDB
	}

	return nil
}

/**
* GetSchema
* @param name string
* @return *Schema
**/
func GetSchema(name string) *Schema {
	list := strs.Split(name, ".")
	if len(list) > 1 {
		db := GetDB(list[0])
		if db == nil {
			return nil
		}

		return db.GetSchema(list[1])
	}

	return nil
}

/**
* GetModel
* @param name string
* @return *Model
**/
func GetModel(name string) *Model {
	var result *Model
	list := strs.Split(name, ".")
	if len(list) > 2 {
		db := GetDB(list[0])
		if db == nil {
			result = nil
		}

		schema := db.GetSchema(list[1])
		if schema == nil {
			result = nil
		}

		result = schema.GetModel(list[2])
	} else if len(list) == 2 {
		db := conn.DefaultDB
		if db == nil {
			result = nil
		}

		schema := db.GetSchema(list[1])
		if schema == nil {
			result = nil
		}

		result = schema.GetModel(list[2])
	} else if len(list) == 1 {
		db := conn.DefaultDB
		if db == nil {
			result = nil
		}

		result = db.GetModel(list[1])
	}

	if err := result.Init(); err != nil {
		console.Logf("jdb", `Error initializing model %s: %v`, name, err)
		return nil
	}

	return result
}

/**
* getTable
* @param name string
* @return *Model
**/
func getTable(table string) *Model {
	var result *Model
	var name string
	list := strs.Split(table, ":")
	if len(list) > 1 {
		name = list[1]
		table = list[0]
	}

	list = strs.Split(table, ".")
	if len(list) == 1 {
		return GetModel(list[0])
	} else if len(list) > 2 {
		if name == "" {
			name = list[2]
		}

		result = GetDB(list[0]).GetSchema(list[1]).GetModel(name)
		if result != nil {
			return result
		}

		schema := GetSchema(list[1])
		if schema == nil {
			return nil
		}

		result = NewModel(schema, name, 1)
		result.Table = list[2]
		result.UseCore = false
	} else if len(list) == 2 {
		if name == "" {
			name = list[1]
		}

		result = GetSchema(list[0]).GetModel(name)
		if result != nil {
			return result
		}

		schema := GetSchema(list[0])
		if schema == nil {
			return nil
		}

		result = NewModel(schema, name, 1)
		result.Table = list[1]
		result.UseCore = false
	}

	if result == nil {
		return nil
	}

	if err := result.Init(); err != nil {
		console.Logf("jdb", `Error initializing model %s: %v`, name, err)
		return nil
	}

	return result
}

/**
* Define
* @param params et.Json
* @return et.Json, error
**/
func define(params et.Json) (et.Json, error) {
	result := et.Json{}
	help := et.Json{
		"message": MSG_INVALID_MODEL_PARAM,
		"help":    "It is required this params.",
		"params": et.Json{
			"name_model": et.Json{
				"schema":    "schema_name",
				"version":   1,
				"integrity": false,
				"fields":    et.Json{},
			},
		},
	}

	for name := range params {
		param := params.Json(name)
		if param.IsEmpty() {
			return et.Json{}, mistake.New(help.ToString())
		}

		schemaName := param.ValStr("jdb", "schema")
		schema := GetSchema(schemaName)
		if schema == nil {
			return et.Json{}, mistake.Newf(MSG_SCHEMA_NOT_FOUND, schemaName)
		}

		version := param.ValInt(1, "version")
		model := NewModel(schema, name, version)
		if model == nil {
			return et.Json{}, mistake.Newf(MSG_MODEL_NOT_FOUND, name)
		}

		integrity := param.ValBool(false, "integrity")
		if integrity {
			model.DefineIntegrity()
		}

		fields := param.Json("fields")
		model.setFields(fields)

		if err := model.Init(); err != nil {
			return et.Json{}, err
		}

		result[name] = model.Describe()
	}

	return result, nil
}

/**
* Describe
* @param kind, name string
* @return et.Json
**/
func describe(kind, name string) (et.Json, error) {
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
		result := GetSchema(name)
		if result == nil {
			return et.Json{}, mistake.Newf(MSG_SCHEMA_NOT_FOUND, name)
		}

		return result.Describe(), nil
	case "model":
		result := GetModel(name)
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

		model := GetModel(list[0])
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
* queryTx
* @param tx *Tx, params et.Json
* @return interface{}, error
**/
func queryTx(tx *Tx, params et.Json) (interface{}, error) {
	from := params.Str("from")
	model := GetModel(from)
	if model == nil {
		return nil, mistake.Newf(MSG_MODEL_NOT_FOUND, from)
	}

	return From(model).
		queryTx(tx, params)
}

/**
* commandsTx
* @param tx *Tx, params et.Json
* @return interface{}, error
**/
func commandsTx(tx *Tx, params et.Json) (et.Items, error) {
	help := et.Json{
		"message": MSG_INVALID_MODEL_PARAM,
		"help":    "It is required this params.",
		"params": et.Json{
			"name_model": et.Json{
				"insert": et.Json{},
				"update": et.Json{},
				"delete": et.Json{},
				"where":  et.Json{},
			},
		},
	}

	result := et.Items{}
	for name := range params {
		model := GetModel(name)
		if model == nil {
			return et.Items{}, mistake.Newf(MSG_MODEL_NOT_FOUND, name)
		}

		param := params.Json(name)
		if param.IsEmpty() {
			return et.Items{}, mistake.New(help.ToString())
		}

		debug := param.ValBool(false, "debug")

		if param["insert"] != nil {
			data := param.Json("insert")
			item, err := model.
				Insert(data).
				setDebug(debug).
				ExecTx(tx)
			if err != nil {
				return et.Items{}, err
			}

			result.AddMany(item.Result)
		}

		if param["update"] != nil {
			data := param.Json("update")
			where := param.Json("where")
			items, err := model.
				Update(data).
				setWheres(where).
				setDebug(debug).
				ExecTx(tx)
			if err != nil {
				return et.Items{}, err
			}

			result.AddMany(items.Result)
		}

		if param["delete"] != nil {
			where := param.Json("where")
			items, err := model.
				delete().
				setWheres(where).
				setDebug(debug).
				ExecTx(tx)
			if err != nil {
				return et.Items{}, err
			}

			result.AddMany(items.Result)
		}

		if param["bulk"] != nil {
			data := param.ArrayJson("bulk")
			items, err := model.
				Bulk(data).
				setDebug(debug).
				ExecTx(tx)
			if err != nil {
				return et.Items{}, err
			}

			result.AddMany(items.Result)
		}

		if param["upsert"] != nil {
			where := param.Json("where")
			data := param.Json("upsert")
			items, err := model.
				Upsert(data).
				setWheres(where).
				setDebug(debug).
				ExecTx(tx)
			if err != nil {
				return et.Items{}, err
			}

			result.AddMany(items.Result)
		}
	}

	return result, nil
}

/**
* ModelDefine
* @param w http.ResponseWriter
* @param r *http.Request
**/
func ModelDefine(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	params := body.Json("params")
	result, err := define(params)
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
	result, err := queryTx(nil, params)
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
	result, err := commandsTx(nil, params)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ANY(w, r, http.StatusOK, result)
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
	result, err := describe(kind, name)
	if err != nil {
		response.JSON(w, r, http.StatusBadRequest, result)
		return
	}

	response.JSON(w, r, http.StatusOK, result)
}

/**
* JSQL
* @param w http.ResponseWriter
* @param r *http.Request
**/
func JSQL(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	result := body

	response.JSON(w, r, http.StatusOK, result)
}
