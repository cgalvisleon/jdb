package sqlite

import (
	"database/sql"

	"github.com/cgalvisleon/et/config"
	"github.com/cgalvisleon/et/envar"
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
	jdb.Register(jdb.SqliteDriver, newDriver, jdb.ConnectParams{
		Id:       config.String("DB_ID", "jdb"),
		Driver:   jdb.SqliteDriver,
		Name:     config.String("DB_NAME", "jdb"),
		UserCore: true,
		Debug:    envar.Bool("DEBUG"),
		Params: et.Json{
			"database": config.String("DB_NAME", "jdb"),
		},
		Validate: []string{
			"DB_NAME",
		},
	})
}
