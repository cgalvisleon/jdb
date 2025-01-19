package postgres

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/jdb"
	_ "github.com/lib/pq"
)

var driver Postgres

type Postgres struct {
	params    et.Json
	connStr   string
	db        *sql.DB
	master    *sql.DB
	connected bool
	version   int
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
