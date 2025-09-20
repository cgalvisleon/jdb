package jdb

import (
	"database/sql"
)

const (
	DriverPostgres = "postgres"
	DriverMysql    = "mysql"
	DriverSqlite   = "sqlite"
	DriverMssql    = "mssql"
	DriverOracle   = "oracle"
)

type Driver interface {
	Connect(db *Database) (*sql.DB, error)
	Load(model *Model) error
	Query(query *Ql) (string, error)
	Command(command *Cmd) (string, error)
}

var drivers map[string]func(db *Database) Driver

func init() {
	drivers = make(map[string]func(db *Database) Driver)
}

func Register(name string, driver func(db *Database) Driver) {
	drivers[name] = driver
}
