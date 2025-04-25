package jdb

import (
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
	Connect(params et.Json) error
	Disconnect() error
	SetMain(params et.Json) error
	// Database
	CreateDatabase(name string) error
	DropDatabase(name string) error
	// User
	GrantPrivileges(username, database string) error
	CreateUser(username, password, confirmation string) error
	ChangePassword(username, password, confirmation string) error
	DeleteUser(username string) error
	// Schema
	LoadSchema(name string) error
	DropSchema(name string) error
	// Model
	LoadModel(model *Model) error
	DropModel(model *Model) error
	MutateModel(model *Model) error
	// Query
	Exec(sql string, arg ...any) error
	Query(sql string, arg ...any) (et.Items, error)
	One(sql string, arg ...any) (et.Item, error)
	Data(source, sql string, arg ...any) (et.Items, error)
	Select(ql *Ql) (et.Items, error)
	Count(ql *Ql) (int, error)
	Exists(ql *Ql) (bool, error)
	// Command
	Command(command *Command) (et.Items, error)
	Sync(command string, data et.Json) error
}

/**
* Register
* @param name string, driver func() Driver
**/
func Register(name string, driver func() Driver) {
	conn.Drivers[name] = driver
}
