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
	SetUser(username, password, confirmation string) error
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
	Query(linq *Linq) (et.Items, error)
	Count(linq *Linq) (int, error)
	Last(linq *Linq) (et.Items, error)
	// Command
	Current(command *Command) (et.Items, error)
	Command(command *Command) (et.Item, error)
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
	drivers[name] = driver
}
