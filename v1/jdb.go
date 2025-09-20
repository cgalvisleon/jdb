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

/**
* Define
* @param definition et.Json
* @return (*jdb.Model, error)
**/
func Define(definition et.Json) (*jdb.Model, error) {
	return jdb.Define(definition)
}

/**
* Select
* @param query et.Json
* @return (*jdb.Ql, error)
**/
func Select(query et.Json) (*jdb.Ql, error) {
	return jdb.Select(query)
}

/**
* Insert
* @param param et.Json
* @return (*jdb.Cmd, error)
**/
func Insert(param et.Json) (*jdb.Cmd, error) {
	return jdb.Insert(param)
}

/**
* Update
* @param param et.Json
* @return (*jdb.Cmd, error)
**/
func Update(param et.Json) (*jdb.Cmd, error) {
	return jdb.Update(param)
}

/**
* Delete
* @param param et.Json
* @return (*jdb.Cmd, error)
**/
func Delete(param et.Json) (*jdb.Cmd, error) {
	return jdb.Delete(param)
}

/**
* Upsert
* @param param et.Json
* @return (*jdb.Cmd, error)
**/
func Upsert(param et.Json) (*jdb.Cmd, error) {
	return jdb.Upsert(param)
}
