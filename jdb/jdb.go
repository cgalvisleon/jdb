package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

/**
* ConnectTo
* @param name, driver string, userCore bool, params Connection
* @return (*DB, error)
**/
func ConnectTo(name, driver string, userCore bool, params et.Json) (*DB, error) {
	return getDatabase(name, driver, userCore, params)
}

/**
* LoadTo
* @param name string
* @return (*DB, error)
**/
func LoadTo(name string) (*DB, error) {
	driver := envar.GetStr("DB_DRIVER", "postgres")
	params := et.Json{
		"database": name,
		"host":     envar.GetStr("DB_HOST", "localhost"),
		"port":     envar.GetInt("DB_PORT", 5432),
		"username": envar.GetStr("DB_USERNAME", "test"),
		"password": envar.GetStr("DB_PASSWORD", "test"),
		"app":      envar.GetStr("DB_APP", "test"),
		"version":  envar.GetInt("DB_VERSION", 15),
	}

	return getDatabase(name, driver, true, params)
}

/**
* Load
* @return (*DB, error)
**/
func Load() (*DB, error) {
	name := envar.GetStr("DB_NAME", "josephine")
	return LoadTo(name)
}

/**
* GetDatabase
* @param name string
* @return (*DB, error)
**/
func GetDatabase(name string) (*DB, error) {
	result, ok := dbs[name]
	if !ok {
		return nil, fmt.Errorf(MSG_DATABASE_NOT_FOUND, name)
	}

	return result, nil
}

/**
* GetModel
* @param database, name string
* @return (*Model, error)
**/
func GetModel(database, name string) (*Model, error) {
	db, ok := dbs[database]
	if !ok {
		return nil, fmt.Errorf("database %s not found", database)
	}

	result, ok := db.Models[name]
	if !ok {
		return nil, fmt.Errorf("model %s not found", name)
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

/**
* From
* @param model *Model
* @return *Ql
**/
func From(model *Model) *Ql {
	db := model.db
	return db.From(model)
}
