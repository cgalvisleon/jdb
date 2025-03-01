package postgres

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
	_ "github.com/lib/pq"
)

type Postgres struct {
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
		params:    et.Json{},
		connected: false,
	}
}

func (s *Postgres) Name() string {
	return jdb.Postgres
}

func init() {
	jdb.Register(jdb.Postgres, NewDriver)
}
