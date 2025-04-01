package postgres

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/msg"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) chain(params et.Json) (string, error) {
	username := params.Str("username")
	password := params.Str("password")
	host := params.Str("host")
	port := params.Int("port")
	database := params.Str("database")
	app := params.Str("app")

	if !utility.ValidStr(username, 0, []string{""}) {
		return "", mistake.Newf(jdb.MSS_PARAM_REQUIRED, "username")
	}

	if !utility.ValidStr(password, 0, []string{""}) {
		return "", mistake.Newf(jdb.MSS_PARAM_REQUIRED, "password")
	}

	if !utility.ValidStr(host, 0, []string{""}) {
		return "", mistake.Newf(jdb.MSS_PARAM_REQUIRED, "host")
	}

	if port == 0 {
		return "", mistake.Newf(jdb.MSS_PARAM_REQUIRED, "port")
	}

	if !utility.ValidStr(database, 0, []string{""}) {
		return "", mistake.Newf(jdb.MSS_PARAM_REQUIRED, "database")
	}

	if !utility.ValidStr(app, 0, []string{""}) {
		return "", mistake.Newf(jdb.MSS_PARAM_REQUIRED, "app")
	}

	driver := s.name
	result := strs.Format(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, driver, username, password, host, port, database, app)

	return result, nil
}

func (s *Postgres) connectTo(connStr string) (*sql.DB, error) {
	db, err := sql.Open(s.name, connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (s *Postgres) connect(params et.Json) error {
	connStr, err := s.chain(params)
	if err != nil {
		return err
	}

	db, err := s.connectTo(connStr)
	if err != nil {
		return err
	}

	s.db = db
	s.params = params
	s.connStr = connStr
	s.connected = s.db != nil
	s.nodeId = params.Int("node_id")
	s.Version()

	return nil
}

func (s *Postgres) connectDefault(params et.Json) error {
	params["database"] = "postgres"
	return s.connect(params)
}

func (s *Postgres) existDatabase(name string) (bool, error) {
	sql := `
	SELECT EXISTS(
	SELECT 1
	FROM pg_database
	WHERE UPPER(datname) = UPPER($1));`
	items, err := s.Query(sql, name)
	if err != nil {
		return false, err
	}

	if items.Count == 0 {
		return false, nil
	}

	return items.Bool(0, "exists"), nil
}

/**
* CreateDatabase
* @param name string
* @return error
**/
func (s *Postgres) CreateDatabase(name string) error {
	if s.db == nil {
		return mistake.Newf(msg.NOT_DRIVER_DB)
	}

	exist, err := s.existDatabase(name)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	sql := jdb.SQLDDL(`CREATE DATABASE $1`, name)
	err = s.Exec(sql)
	if err != nil {
		return err
	}

	console.Logf(s.name, `Database %s created`, name)

	return nil
}

/**
* DropDatabase
* @param name string
* @return error
**/
func (s *Postgres) DropDatabase(name string) error {
	if s.db == nil {
		return mistake.Newf(msg.NOT_DRIVER_DB)
	}

	exist, err := s.existDatabase(name)
	if err != nil {
		return err
	}

	if !exist {
		return nil
	}

	sql := jdb.SQLDDL(`DROP DATABASE $1`, name)
	err = s.Exec(sql)
	if err != nil {
		return err
	}

	console.Logf(s.name, `Database %s droped`, name)

	return nil
}

/**
* Connect
* @param params et.Json
* @return error
**/
func (s *Postgres) Connect(params et.Json) error {
	database := params.Str("database")
	err := s.connectDefault(params)
	if err != nil {
		return err
	}

	err = s.CreateDatabase(database)
	if err != nil {
		return err
	}

	params["database"] = database
	err = s.connect(params)
	if err != nil {
		return err
	}

	console.Logf(s.name, `Connected to %s:%s`, params.Str("host"), database)

	return nil
}

/**
* Disconnect
* @return error
**/
func (s *Postgres) Disconnect() error {
	if !s.connected {
		return nil
	}

	if s.db != nil {
		s.db.Close()
	}

	return nil
}

/**
* Version
* @return int
**/
func (s *Postgres) Version() int {
	if !s.connected {
		return 0
	}

	if s.db == nil {
		return 0
	}

	if s.version != 0 {
		return s.version
	}

	var version string
	err := s.db.QueryRow("SELECT version();").Scan(&version)
	if err != nil {
		return 0
	}

	split := strings.Split(version, ".")
	v, err := strconv.Atoi(split[0])
	if err != nil {
		v = 13
	}

	s.version = v

	return s.version
}
