package jdb

import (
	"github.com/cgalvisleon/et/et"
)

var (
	Postgres = "postgres"
	Josefine = "josefine"
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
	CreateUser(username, password, confirmation string) error
	ChangePassword(username, password, confirmation string) error
	DeleteUser(username string) error
	// Schema
	CreateSchema(name string) error
	DropSchema(name string) error
	// Model
	LoadModel(model *Model) error
	// Query
	Exec(sql string, params ...any) error
	SQL(sql string, params ...any) (et.Items, error)
	One(sql string, params ...any) (et.Item, error)
	Query(ql *Ql) (et.Items, error)
	Count(ql *Ql) (int, error)
	// Command
	Command(command *Command) (et.Items, error)
	// Series
	GetSerie(tag string) int64
	NextCode(tag, prefix string) string
	SetSerie(tag string, val int) int64
	CurrentSerie(tag string) int64
	// Key Value
	SetKey(key string, value []byte) error
	GetKey(key string) (et.KeyValue, error)
	DeleteKey(key string) error
	FindKeys(search string, page, rows int) (et.List, error)
}

/**
* SetDriver
**/
func Register(name string, driver func() Driver) {
	drivers[name] = driver
}
