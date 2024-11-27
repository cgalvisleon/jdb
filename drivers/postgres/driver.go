package postgres

import (
	"database/sql"

	jdb "github.com/cgalvisl/jdb/pkg"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/ws"
	_ "github.com/lib/pq"
)

const DriverName = "postgres"

var driver Postgres

type Postgres struct {
	db        *sql.DB
	master    *sql.DB
	ws        *ws.Hub
	params    et.Json
	connStr   string
	connected bool
}

func NewDriver() jdb.Driver {
	return &Postgres{}
}

func (s *Postgres) Name() string {
	return DriverName
}

func (s *Postgres) Connect(params et.Json) error {
	if params.Str("username") == "" {
		return logs.Alertf(jdb.MSS_PARAM_REQUIRED, "username")
	}

	if params.Str("password") == "" {
		return logs.Alertf(jdb.MSS_PARAM_REQUIRED, "password")
	}

	if params.Str("host") == "" {
		return logs.Alertf(jdb.MSS_PARAM_REQUIRED, "host")
	}

	if params.Int("port") == 0 {
		return logs.Alertf(jdb.MSS_PARAM_REQUIRED, "port")
	}

	if params.Str("database") == "" {
		return logs.Alertf(jdb.MSS_PARAM_REQUIRED, "database")
	}

	if params.Str("app") == "" {
		return logs.Alertf(jdb.MSS_PARAM_REQUIRED, "app")
	}

	driver := DriverName
	user := params.Str("user")
	password := params.Str("password")
	host := params.Str("host")
	port := params.Int("port")
	database := params.Str("database")
	app := params.Str("app")
	connStr := strs.Format(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, driver, user, password, host, port, database, app)

	db, err := sql.Open(DriverName, connStr)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.params = params
	s.connStr = connStr
	s.db = db
	s.connected = true

	return nil
}

func (s *Postgres) Disconnect() error {
	if !s.connected {
		return nil
	}

	if s.db != nil {
		s.db.Close()
	}

	return nil
}

func init() {
	jdb.Register(DriverName, NewDriver)
}
