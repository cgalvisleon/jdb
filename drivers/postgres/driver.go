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
* @return error
**/
func (s *Postgres) Load(model *jdb.Model) error {
	model.Table = fmt.Sprintf("%s.%s", model.Schema, model.Name)
	model.SetInit()

	return nil
}

/**
* Query
* @param query *jdb.Ql
* @return (string, error)
**/
func (s *Postgres) Query(query *jdb.Ql) (string, error) {
	defintion := query.ToJson()
	sql, err := s.buildQuery(defintion)
	if err != nil {
		return "", err
	}

	query.SQL = sql
	console.Debug("sql:\n\t", sql)

	return query.SQL, nil
}

/**
* Cmd
* @param command *jdb.Cmd
* @return (string, error)
**/
func (s *Postgres) Command(command *jdb.Cmd) (string, error) {
	return command.SQL, nil
}
