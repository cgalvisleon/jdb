package oci

import (
	"database/sql"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
	_ "github.com/lib/pq"
)

type Postgres struct {
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
	return &Postgres{
		name:      jdb.PostgresDriver,
		params:    et.Json{},
		connected: false,
	}
}

func (s *Postgres) Name() string {
	return s.name
}

func init() {
	jdb.Register(jdb.PostgresDriver, NewDriver)
	envar.UpSetStr("DB_DRIVER", jdb.PostgresDriver)
}
