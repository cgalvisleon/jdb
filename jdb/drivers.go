package jdb

import (
	"github.com/cgalvisleon/et/et"
)

const (
	SqliteDriver   = "sqlite"
	PostgresDriver = "postgres"
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
	// Core
	CreateCore() error
	// User
	GrantPrivileges(username, database string) error
	CreateUser(username, password, confirmation string) error
	ChangePassword(username, password, confirmation string) error
	DeleteUser(username string) error
	// Schema
	CreateSchema(name string) error
	DropSchema(name string) error
	// Model
	LoadTable(model *Model) (bool, error)
	CreateModel(model *Model) error
	DropModel(model *Model) error
	SaveModel(model *Model) error
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
	// Series
	GetSerie(tag string) int64
	NextCode(tag, prefix string) string
	SetSerie(tag string, val int) int64
	CurrentSerie(tag string) int64
}

/**
* SetDriver
**/
func Register(name string, driver func() Driver) {
	Jdb.Drivers[name] = driver
}
