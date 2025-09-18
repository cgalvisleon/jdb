package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
)

/**
* setModel
* @param id string, data et.Json
* @return error
**/
func setModel(id string, data et.Json, debug bool) error {
	if debug {
		console.Debugf("%s:%s", id, data.ToString())
	}

	return nil
}
