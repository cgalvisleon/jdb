package postgres

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
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

func (s *Postgres) chain(params et.Json) (string, error) {
	username := params.Str("username")
	password := params.Str("password")
	host := params.Str("host")
	port := params.Int("port")
	database := params.Str("database")
	app := params.Str("app")

	if !utility.ValidStr(username, 0, []string{""}) {
		return "", logs.Alertf(jdb.MSS_PARAM_REQUIRED, "username")
	}

	if !utility.ValidStr(password, 0, []string{""}) {
		return "", logs.Alertf(jdb.MSS_PARAM_REQUIRED, "password")
	}

	if !utility.ValidStr(host, 0, []string{""}) {
		return "", logs.Alertf(jdb.MSS_PARAM_REQUIRED, "host")
	}

	if port == 0 {
		return "", logs.Alertf(jdb.MSS_PARAM_REQUIRED, "port")
	}

	if !utility.ValidStr(database, 0, []string{""}) {
		return "", logs.Alertf(jdb.MSS_PARAM_REQUIRED, "database")
	}

	if !utility.ValidStr(app, 0, []string{""}) {
		return "", logs.Alertf(jdb.MSS_PARAM_REQUIRED, "app")
	}

	driver := jdb.Postgres
	result := strs.Format(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, driver, username, password, host, port, database, app)

	return result, nil
}

func (s *Postgres) connectTo(connStr string) (*sql.DB, error) {
	db, err := sql.Open(jdb.Postgres, connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (s *Postgres) Connect(params et.Json) error {
	connStr, err := s.chain(params)
	if err != nil {
		return err
	}

	db, err := s.connectTo(connStr)
	if err != nil {
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

func (s *Postgres) SetMain(params et.Json) error {

	return nil
}

func init() {
	jdb.Register(jdb.Postgres, NewDriver)
}
