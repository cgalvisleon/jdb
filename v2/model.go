package jdb

import (
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/congo"
)

/**
* DefineModel
* @param definition et.Json
* @return (*jdb.Model, error)
**/
func DefineModel(definition et.Json) (*jdb.Model, error) {
	return jdb.DefineModel(definition)
}
