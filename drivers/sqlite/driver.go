package sqlite

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
	_ "modernc.org/sqlite"
)

type SqlLite struct {
	name      string
	params    et.Json
	connStr   string
	db        *sql.DB
	connected bool
	version   int
	nodeId    int
}

func newDriver() jdb.Driver {
	return &SqlLite{
		name:      jdb.SqliteDriver,
		params:    et.Json{},
		connected: false,
	}
}

func (s *SqlLite) Name() string {
	return s.name
}

func init() {
	jdb.Register(jdb.SqliteDriver, newDriver)
}
