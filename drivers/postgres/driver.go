package postgres

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
	_ "github.com/lib/pq"
)

type Postgres struct {
	name      string
	params    et.Json
	connStr   string
	db        *sql.DB
	connected bool
	version   int
	nodeId    int
}

func newDriver() jdb.Driver {
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
	jdb.Register(jdb.PostgresDriver, newDriver)
}
