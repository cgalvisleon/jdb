package postgres

import (
	"database/sql"
	"fmt"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/jdb/jdb"
)

var driver = "postgres"

func init() {
	jdb.Register(driver, newDriver)
}

type Postgres struct {
	name       string      `json:"-"`
	database   *jdb.DB     `json:"-"`
	connection *Connection `json:"-"`
}

/**
* newDriver
* @param database *jdb.DB
* @return jdb.Driver
**/
func newDriver(database *jdb.DB) jdb.Driver {
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

	result.connection.Load(result.database.Connection)

	return result
}

/**
* Connect
* @param connection jdb.ConnectParams
* @return *sql.DB, error
**/
func (s *Postgres) Connect(database *jdb.DB) (*sql.DB, error) {
	s.database = database
	s.name = database.Name

	defaultChain, err := s.connection.defaultChain()
	if err != nil {
		return nil, err
	}

	db, err := connectTo(defaultChain)
	if err != nil {
		return nil, err
	}

	err = s.CreateDatabase(db, database.Name)
	if err != nil {
		return nil, err
	}

	if db != nil {
		err := db.Close()
		if err != nil {
			return nil, err
		}
	}

	s.connection.Database = database.Name
	chain, err := s.connection.chain()
	if err != nil {
		return nil, err
	}

	db, err = connectTo(chain)
	if err != nil {
		return nil, err
	}

	if database.UseCore {
		if err := triggerRecords(db); err != nil {
			return nil, err
		}
	}

	logs.Logf(driver, `Connected to %s:%s`, s.connection.Host, s.connection.Database)

	return db, nil
}

/**
* Load
* @param model *Model
* @return (string, error)
**/
func (s *Postgres) Load(model *jdb.Model) (string, error) {
	model.Table = fmt.Sprintf("%s.%s", model.Schema, model.Name)
	exists, err := existTable(s.database.Db, model.Schema, model.Name)
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
	if ql.IsDebug {
		logs.Debug("query:\n", sql)
	}

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
	if cmd.IsDebug {
		logs.Debug("command:\n", sql)
	}

	return cmd.SQL, nil
}
