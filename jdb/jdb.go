package jdb

import (
	"slices"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/strs"
)

type JDB struct {
	Drivers map[string]func() Driver
	DBS     []*DB
	Version string
}

var Jdb *JDB

func init() {
	Jdb = &JDB{
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
	drivers := []string{}
	for key := range s.Drivers {
		drivers = append(drivers, key)
	}
	dbs := []et.Json{}
	for _, db := range s.DBS {
		dbs = append(dbs, db.Describe())
	}

	return et.Json{
		"drivers": drivers,
		"dbs":     dbs,
		"version": s.Version,
	}
}

/**
* ConnectTo
* @param params et.Json
* @return *DB, error
**/
func ConnectTo(params ConnectParams) (*DB, error) {
	if Jdb == nil {
		return nil, mistake.New(MSG_JDB_NOT_DEFINED)
	}

	name := params.Name
	if name == "" {
		return nil, mistake.New(MSG_DATABASE_NOT_DEFINED)
	}

	driver := params.Driver
	if driver == "" {
		return nil, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	result, err := NewDatabase(name, driver)
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
			case "StatusField":
				StatusField = ColumnField(value)
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
	name = Name(name)
	idx := slices.IndexFunc(Jdb.DBS, func(db *DB) bool { return db.Name == name })
	if idx != -1 {
		return Jdb.DBS[idx]
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
	if len(Jdb.DBS) == 0 {
		return nil
	}

	name = Name(name)
	var db *DB
	var result *Schema
	var err error
	list := strs.Split(name, ".")
	switch len(list) {
	case 1:
		db = Jdb.DBS[0]
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
		result, err = NewSchema(db, name)
		if err != nil {
			return nil
		}
	}

	return result
}

/**
* GetModel
* @param name string, isCreate bool
* @return *Model
**/
func GetModel(name string, isCreate bool) *Model {
	if len(Jdb.DBS) == 0 {
		return nil
	}

	var db *DB
	var schema *Schema
	var result *Model
	list := strs.Split(name, ".")
	switch len(list) {
	case 1: // model
		db = Jdb.DBS[0]
		if db == nil {
			return nil
		}

		schema = db.schemas[0]
		if schema == nil {
			return nil
		}

		result = schema.GetModel(name)
		if result != nil {
			return result
		}
	case 2: // schema, model
		db = Jdb.DBS[0]
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
* GetField
* @param name string
* @return *Field
**/
func GetField(name string) *Field {
	list := strs.Split(name, ".")

	switch len(list) {
	case 2: // model
		modelName := list[0]
		model := GetModel(modelName, false)
		if model == nil {
			return nil
		}

		return model.GetField(list[1])
	case 3: // schema, model
		modelName := strs.Format(`%s.%s`, list[0], list[1])
		model := GetModel(modelName, false)
		if model == nil {
			return nil
		}

		return model.GetField(list[2])
	case 4: // db, schema, model
		modelName := strs.Format(`%s.%s.%s`, list[0], list[1], list[2])
		model := GetModel(modelName, false)
		if model == nil {
			return nil
		}

		return model.GetField(list[3])
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

		return Jdb.Describe(), nil
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

		field := model.GetField(list[3])
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
		model := GetModel(insert, false)
		if model == nil {
			return nil, mistake.Newf(MSG_MODEL_NOT_FOUND, insert)
		}

		conn = model.Insert(data)
	} else if update != "" {
		model := GetModel(update, false)
		if model == nil {
			return nil, mistake.Newf(MSG_MODEL_NOT_FOUND, update)
		}

		conn = model.Update(data).
			setWhere(where)
	} else if delete != "" {
		model := GetModel(delete, false)
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
