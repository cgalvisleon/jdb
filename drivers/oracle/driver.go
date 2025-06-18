package oracle

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
	_ "github.com/lib/pq"
)

type Oracle struct {
	name      string
	params    et.Json
	connStr   string
	db        *sql.DB
	connected bool
	version   int
	nodeId    int
}

func newDriver() jdb.Driver {
	return &Oracle{
		name:      jdb.PostgresDriver,
		params:    et.Json{},
		connected: false,
	}
}

func (s *Oracle) Name() string {
	return s.name
}

func init() {
	jdb.Register(jdb.PostgresDriver, newDriver)
}
