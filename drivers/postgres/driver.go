package postgres

import (
	"database/sql"

	"github.com/cgalvisleon/et/config"
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
	connected bool
	version   int
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
	jdb.Register(jdb.PostgresDriver, newDriver, jdb.ConnectParams{
		Driver:   jdb.PostgresDriver,
		Name:     config.String("DB_NAME", "jdb"),
		UserCore: true,
		Debug:    envar.Bool("DEBUG"),
		Params: et.Json{
			"database": config.String("DB_NAME", "jdb"),
			"host":     config.String("DB_HOST", "localhost"),
			"port":     config.Int("DB_PORT", 5432),
			"username": config.String("DB_USER", "admin"),
			"password": config.String("DB_PASSWORD", "admin"),
			"app":      config.App.Name,
		},
		Validate: []string{
			"DB_NAME",
			"DB_HOST",
			"DB_PORT",
			"DB_USER",
			"DB_PASSWORD",
		},
	})
}
