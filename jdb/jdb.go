package jdb

import (
	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
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
