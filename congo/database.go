package jdb

import "github.com/cgalvisleon/et/et"

var dbs map[string]*Database

func init() {
	dbs = make(map[string]*Database)
}

type Database struct {
	Name   string            `json:"name"`
	Models map[string]*Model `json:"models"`
	Driver Driver            `json:"driver"`
}

func (s *Database) ToJson() et.Json {
	return et.Json{
		"name":   s.Name,
		"models": s.Models,
	}
}

func getDatabase(name string) *Database {
	result, ok := dbs[name]
	if !ok {
		result = &Database{
			Name:   name,
			Models: make(map[string]*Model),
			Driver: drivers["postgres"],
		}
		dbs[name] = result
	}

	return result
}

func (s *Database) loadModel(model *Model) (string, error) {
	s.Models[model.Id] = model
	return s.Driver.Load(model)
}
