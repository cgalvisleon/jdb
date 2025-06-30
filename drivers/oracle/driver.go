package oracle

import (
	"database/sql"

	"github.com/cgalvisleon/et/config"
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
}

func newDriver() jdb.Driver {
	return &Oracle{
		name:      jdb.OracleDriver,
		params:    et.Json{},
		connected: false,
	}
}

func (s *Oracle) Name() string {
	return s.name
}

func init() {
	jdb.Register(jdb.OracleDriver, newDriver, jdb.ConnectParams{
		Id:     config.String("DB_ID", "jdb"),
		Driver: jdb.OracleDriver,
		Name:   config.String("DB_NAME", "jdb"),
		Params: et.Json{
			"database":     config.String("DB_NAME", "jdb"),
			"host":         config.String("DB_HOST", "localhost"),
			"port":         config.Int("DB_PORT", 5432),
			"username":     config.String("DB_USER", "admin"),
			"password":     config.String("DB_PASSWORD", "admin"),
			"app":          config.App.Name,
			"service_name": config.String("ORA_DB_SERVICE_NAME_ORACLE", "jdb"),
			"ssl":          config.Bool("ORA_DB_SSL_ORACLE", false),
			"ssl_verify":   config.Bool("ORA_DB_SSL_VERIFY_ORACLE", false),
			"version":      config.Int("ORA_DB_VERSION_ORACLE", 19),
		},
		UserCore: true,
		Validate: []string{
			"DB_NAME",
			"DB_HOST",
			"DB_PORT",
			"DB_USER",
			"DB_PASSWORD",
			"ORA_DB_SERVICE_NAME_ORACLE",
			"ORA_DB_SSL_ORACLE",
			"ORA_DB_SSL_VERIFY_ORACLE",
			"ORA_DB_VERSION_ORACLE",
		},
	})
}
