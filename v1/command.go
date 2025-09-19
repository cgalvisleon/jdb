package jdb

import (
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/congo"
)

/**
* Insert
* @param param et.Json
* @return (*jdb.Command, error)
**/
func Insert(param et.Json) (*jdb.Command, error) {
	return jdb.Insert(param)
}

/**
* Update
* @param param et.Json
* @return (*jdb.Command, error)
**/
func Update(param et.Json) (*jdb.Command, error) {
	return jdb.Update(param)
}

/**
* Delete
* @param param et.Json
* @return (*jdb.Command, error)
**/
func Delete(param et.Json) (*jdb.Command, error) {
	return jdb.Delete(param)
}

/**
* Upsert
* @param param et.Json
* @return (*jdb.Command, error)
**/
func Upsert(param et.Json) (*jdb.Command, error) {
	return jdb.Upsert(param)
}
