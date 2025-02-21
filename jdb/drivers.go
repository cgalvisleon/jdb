package jdb

import (
	"time"

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
	GrantPrivileges(username, dbName string) error
	CreateUser(username, password, confirmation string) error
	ChangePassword(username, password, confirmation string) error
	DeleteUser(username string) error
	// Schema
	CreateSchema(name string) error
	DropSchema(name string) error
	// Model
	LoadTable(model *Model) (bool, error)
	LoadModel(model *Model) error
	DropModel(model *Model) error
	// Query
	Exec(sql string, params ...any) error
	Query(sql string, params ...any) (et.Items, error)
	Data(source, sql string, params ...any) (et.Items, error)
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
	// Key Value
	SetKey(key string, value []byte) error
	GetKey(key string) (et.KeyValue, error)
	DeleteKey(key string) error
	FindKeys(search string, page, rows int) (et.List, error)
	// Flow
	SetFlow(name string, value []byte) error
	GetFlow(id string) (Flow, error)
	DeleteFlow(id string) error
	FindFlows(search string, page, rows int) (et.List, error)
	// Cache
	SetCache(key string, value []byte, duration time.Duration) error
	GetCache(key string) (et.KeyValue, error)
	DeleteCache(key string) error
	CleanCache() error
	FindCache(search string, page, rows int) (et.List, error)
}

/**
* SetDriver
**/
func Register(name string, driver func() Driver) {
	Jdb.Drivers[name] = driver
}
