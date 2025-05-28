package jdb

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
)

const (
	PostgresDriver = "postgres"
	SqliteDriver   = "sqlite"
	MysqlDriver    = "mysql"
	OracleDriver   = "oracle"
)

type Driver interface {
	Name() string
	Connect(params et.Json) (*sql.DB, error)
	Disconnect() error
	/* Model */
	LoadModel(model *Model) error
	DropModel(model *Model) error
	MutateModel(model *Model) error
	/* Ql */
	Select(ql *Ql) (et.Items, error)
	Count(ql *Ql) (int, error)
	Exists(ql *Ql) (bool, error)
	/* Command */
	Command(command *Command) (et.Items, error)
}

/**
* Register
* @param name string, driver func() Driver
**/
func Register(name string, driver func() Driver) {
	conn.Drivers[name] = driver
}
