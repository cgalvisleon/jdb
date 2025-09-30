package postgres

import (
	"fmt"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/jdb/jdb"
)

var driver = "postgres"

func init() {
	jdb.Register(driver, newDriver)
}

type Postgres struct {
	name       string        `json:"-"`
	database   *jdb.Database `json:"-"`
	connection *Connection   `json:"-"`
}

/**
* newDriver
* @param database *jdb.Database
* @return jdb.Driver
**/
func newDriver(database *jdb.Database) jdb.Driver {
	result := &Postgres{
		database: database,
		name:     database.Name,
		connection: &Connection{
			Database: envar.GetStr("DB_NAME", "jdb"),
			Host:     envar.GetStr("DB_HOST", "localhost"),
			Port:     envar.GetInt("DB_PORT", 5432),
			Username: envar.GetStr("DB_USER", "admin"),
			Password: envar.GetStr("DB_PASSWORD", "admin"),
			App:      envar.GetStr("APP_NAME", "jdb"),
			Version:  envar.GetInt("DB_VERSION", 13),
		},
	}

	result.connection.load(result.database.Connection)

	return result
}

/**
* Load
* @param model *Model
* @return (string, error)
**/
func (s *Postgres) Load(model *jdb.Model) (string, error) {
	model.Table = fmt.Sprintf("%s.%s", model.Schema, model.Name)
	exists, err := s.existTable(s.database.Db, model.Schema, model.Name)
	if err != nil {
		return "", err
	}

	if exists {
		model.SetInit()
		return "", nil
	}

	sql, err := s.buildModel(model)
	if err != nil {
		return "", err
	}

	return sql, nil
}

/**
* Query
* @param ql *jdb.Ql
* @return (string, error)
**/
func (s *Postgres) Query(ql *jdb.Ql) (string, error) {
	sql, err := s.buildQuery(ql)
	if err != nil {
		return "", err
	}

	ql.SQL = sql
	console.Debug("query:\n", sql)

	return ql.SQL, nil
}

/**
* Cmd
* @param command *jdb.Cmd
* @return (string, error)
**/
func (s *Postgres) Command(cmd *jdb.Cmd) (string, error) {
	sql, err := s.buildCommand(cmd)
	if err != nil {
		return "", err
	}

	cmd.SQL = sql
	console.Debug("command:\n", sql)

	return cmd.SQL, nil
}
