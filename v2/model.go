package jdb

import (
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/congo"
)

/**
* NewModel
* @param definition et.Json
* @return (*jdb.Model, error)
**/
func NewModel(definition et.Json) (*jdb.Model, error) {
	return jdb.DefineModel(definition)
}
