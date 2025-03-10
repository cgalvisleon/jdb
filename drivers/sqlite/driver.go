package sqlite

import (
	"database/sql"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/jdb/jdb"
	_ "modernc.org/sqlite"
)

type SqlLite struct {
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
	jdb.Register(jdb.SqliteDriver, NewDriver)
	envar.UpSetStr("DB_DRIVER", jdb.SqliteDriver)
}

func (s *SqlLite) Disconnect() error {
	return mistake.New(MSG_FUNCION_NOT_FOUND)
}

func (s *SqlLite) SetMain(arg et.Json) error {
	return mistake.New(MSG_FUNCION_NOT_FOUND)
}
