package jdb

import (
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

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
