package oci

import (
	"database/sql"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
	_ "github.com/lib/pq"
)

type Oracle struct {
	name      string
	params    et.Json
	connStr   string
	db        *sql.DB
	master    *sql.DB
	connected bool
	version   int
	nodeId    int
}

func NewDriver() jdb.Driver {
	return &Oracle{
		name:      jdb.OracleDriver,
		params:    et.Json{},
		connected: false,
	}
}

func (s *Oracle) Name() string {
	return s.name
}

func init() {
	jdb.Register(jdb.OracleDriver, NewDriver)
	envar.UpSetStr("DB_DRIVER", jdb.OracleDriver)
}
