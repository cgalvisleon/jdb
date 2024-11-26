package jdb

var (
	drivers map[string]func() Driver
	dbs     map[string]*Database
	schemas map[string]*Schema
	models  map[string]*Model
)

func init() {
	drivers = map[string]func() Driver{}
	dbs = map[string]*Database{}
	schemas = map[string]*Schema{}
	models = map[string]*Model{}
}
