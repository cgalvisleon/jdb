package mysql

import (
	"database/sql"

	"github.com/cgalvisleon/et/config"
	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	name      string
	params    et.Json
	connStr   string
	db        *sql.DB
	connected bool
	version   int
}

func newDriver() jdb.Driver {
	return &Mysql{
		name:      jdb.MysqlDriver,
		params:    et.Json{},
		connected: false,
	}
}

func (s *Mysql) Name() string {
	return s.name
}

func init() {
	jdb.Register(jdb.PostgresDriver, newDriver, jdb.ConnectParams{
		Driver:   jdb.OracleDriver,
		Name:     config.String("DB_NAME", "jdb"),
		UserCore: true,
		Debug:    envar.Bool("DEBUG"),
		Params: et.Json{
			"database": config.String("DB_NAME", "jdb"),
			"host":     config.String("DB_HOST", "localhost"),
			"port":     config.Int("DB_PORT", 3306),
			"username": config.String("DB_USER", "admin"),
			"password": config.String("DB_PASSWORD", "admin"),
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
