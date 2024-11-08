package jdb

import (
	"github.com/cgalvisleon/et/et"
)

type Driver interface {
	Name() string
	Connect(params et.Json) error
	Disconnect() error
	// Database
	CreateDatabase(name string) error
	DropDatabase(name string) error
	// User
	SetUser(username, password, confirmation string) error
	DeleteUser(username string) error
	SetParams(data et.Json) error
	// Schema
	CreateSchema(name string) error
	DropSchema(name string) error
	// Model
	CreateModel(model *Model) error
	// Query
	SQL(sql string, params ...interface{}) (et.Items, error)
	Query(linq *Linq) (et.Items, error)
	Count(linq *Linq) (int, error)
	Last(linq *Linq) (et.Items, error)
	// Command
	Current(command *Command) (et.Items, error)
	Command(command *Command) (et.Item, error)
	// Series
	GetIndex(tag string) int64
	SetIndex(tag string, val int) int64
}

var Drivers map[string]*Driver

/**
* SetDriver
**/
func SetDriver(d *Driver) {
	Drivers[(*d).Name()] = d
}

func init() {
	Drivers = map[string]*Driver{}
}
