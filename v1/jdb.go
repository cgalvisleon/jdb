package jdb

import (
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* ConnectTo
* @param name, driver string, params et.Json
* @return (*jdb.Database, error)
**/
func ConnectTo(name, driver string, params et.Json) (*jdb.Database, error) {
	return jdb.ConnectTo(name, driver, params)
}

/**
* LoadTo
* @param name string
* @return (*jdb.Database, error)
**/
func LoadTo(name string) (*jdb.Database, error) {
	return jdb.LoadTo(name)
}

/**
* Load
* @return (*jdb.Database, error)
**/
func Load() (*jdb.Database, error) {
	return jdb.Load()
}
