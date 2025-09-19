package jdb

import (
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/congo"
)

/**
* Query
* @param query et.Json
* @return (*jdb.Ql, error)
**/
func Query(query et.Json) (*jdb.Ql, error) {
	return jdb.Query(query)
}
