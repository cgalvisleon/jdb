package postgres

import (
	"database/sql"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/msg"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) ExistDatabase(name string) (bool, error) {
	sql := `
	SELECT EXISTS(
	SELECT 1
	FROM pg_database
	WHERE UPPER(datname) = UPPER($1));`
	items, err := s.SQL(sql, name)
	if err != nil {
		return false, err
	}

	return items.Ok, nil
}

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

	return nil
}

func (s *Postgres) CreateDatabase(name string) error {
	if s.db == nil {
		return mistake.Newf(msg.NOT_DRIVER_DB)
	}

	exist, err := s.ExistDatabase(name)
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

	console.Logf(jdb.Postgres, `Database %s created`, name)

	return nil
}

func (s *Postgres) DropDatabase(name string) error {
	if s.db == nil {
		return mistake.Newf(msg.NOT_DRIVER_DB)
	}

	exist, err := s.ExistDatabase(name)
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

	console.Logf(jdb.Postgres, `Database %s droped`, name)

	return nil
}

func (s *Postgres) Connect(params et.Json) error {
	database := params.Str("database")
	params["database"] = "postgres"
	err := s.connect(params)
	if err != nil {
		return err
	}

	err = s.CreateDatabase(database)
	if err != nil {
		return err
	}

	s.params["database"] = database
	err = s.connect(s.params)
	if err != nil {
		return err
	}

	console.Logf(jdb.Postgres, `Connected to %s:%s`, params.Str("host"), database)

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
