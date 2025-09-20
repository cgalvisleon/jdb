package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

/**
* ConnectTo
* @param name, driver string, params et.Json
* @return (*Database, error)
**/
func ConnectTo(name, driver string, params et.Json) (*Database, error) {
	return getDatabase(name, driver, params)
}

/**
* LoadTo
* @param name string
* @return (*Database, error)
**/
func LoadTo(name string) (*Database, error) {
	driver := envar.GetStr("DB_DRIVER", "postgres")
	return getDatabase(name, driver, et.Json{
		"database": name,
		"host":     envar.GetStr("DB_HOST", "localhost"),
		"port":     envar.GetInt("DB_PORT", 5432),
		"username": envar.GetStr("DB_USERNAME", "test"),
		"password": envar.GetStr("DB_PASSWORD", "test"),
		"app":      envar.GetStr("DB_APP", "test"),
		"version":  envar.GetInt("DB_VERSION", 15),
	})
}

/**
* Load
* @return (*Database, error)
**/
func Load() (*Database, error) {
	name := envar.GetStr("DB_NAME", "josephine")
	return LoadTo(name)
}

/**
* GetDatabase
* @param name string
* @return (*Database, error)
**/
func GetDatabase(name string) (*Database, error) {
	result, ok := dbs[name]
	if !ok {
		return nil, fmt.Errorf(MSG_DATABASE_NOT_FOUND, name)
	}

	return result, nil
}

/**
* GetModel
* @param database, schema, name string
* @return (*Model, error)
**/
func GetModel(database, schema, name string) (*Model, error) {
	db, ok := dbs[database]
	if !ok {
		return nil, fmt.Errorf("database %s not found", database)
	}

	id := fmt.Sprintf("%s.%s", schema, name)
	result, ok := db.Models[id]
	if !ok {
		return nil, fmt.Errorf("model %s not found", id)
	}

	return result, nil
}

/**
* Select
* @param query et.Json
* @return (*Ql, error)
**/
func Select(query et.Json) (*Ql, error) {
	database := query.String("database")
	if !utility.ValidStr(database, 0, []string{}) {
		return nil, fmt.Errorf(MSG_DATABASE_REQUIRED)
	}

	db, err := GetDatabase(database)
	if err != nil {
		return nil, err
	}

	return db.Select(query)
}
