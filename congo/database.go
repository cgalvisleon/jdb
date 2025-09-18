package jdb

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
)

var dbs map[string]*Database

func init() {
	dbs = make(map[string]*Database)
}

type Database struct {
	Name       string            `json:"name"`
	Models     map[string]*Model `json:"models"`
	UseCore    bool              `json:"use_core"`
	Connection et.Json           `json:"-"`
	driver     Driver            `json:"-"`
	db         *sql.DB           `json:"-"`
}

/**
* ToJson
* @return et.Json
**/
func (s *Database) ToJson() et.Json {
	bt, err := json.Marshal(s)
	if err != nil {
		return et.Json{}
	}

	var result et.Json
	err = json.Unmarshal(bt, &result)
	if err != nil {
		return et.Json{}
	}

	return result
}

/**
* GetDatabase
* @param name, driver string, params et.Json
* @return (*Database, error)
**/
func getDatabase(name, driver string, params et.Json) (*Database, error) {
	result, ok := dbs[name]
	if !ok {
		if _, ok := drivers[driver]; !ok {
			return nil, fmt.Errorf(MSG_DRIVER_NOT_FOUND, driver)
		}

		result = &Database{
			Name:       name,
			Models:     make(map[string]*Model),
			Connection: params,
		}
		result.driver = drivers[driver](result)
		err := result.load()
		if err != nil {
			return nil, err
		}

		dbs[name] = result
	}

	return result, nil
}

/**
* load
* @return error
**/
func (s *Database) load() error {
	if s.driver == nil {
		return fmt.Errorf(MSG_DRIVER_REQUIRED)
	}

	db, err := s.driver.Connect(s)
	if err != nil {
		return err
	}

	s.db = db

	if s.UseCore {
		err := initCore()
		if err != nil {
			console.Panic(err)
		}
	}

	return nil
}

/**
* initModel
* @param model *Model
* @return error
**/
func (s *Database) init(model *Model) error {
	err := s.driver.Load(model)
	if err != nil {
		return err
	}

	err = model.save()
	if err != nil {
		return err
	}

	return nil
}

/**
* query
* @param query *Ql
* @return (et.Items, error)
**/
func (s *Database) query(query *Ql) (et.Items, error) {
	if err := query.validate(); err != nil {
		return et.Items{}, err
	}

	result, err := s.driver.Query(query)
	if err != nil {
		return et.Items{}, err
	}

	if query.isDebug {
		console.Debugf("query:%s", query.toJson().ToString())
	}

	return result, nil
}

/**
* exists
* @param query *Ql
* @return (bool, error)
**/
func (s *Database) exists(query *Ql) (bool, error) {
	if err := query.validate(); err != nil {
		return false, err
	}

	result, err := s.driver.Exists(query)
	if err != nil {
		return false, err
	}

	if query.isDebug {
		console.Debugf("exists:%s", query.toJson().ToString())
	}

	return result, nil
}

/**
* command
* @param command *Command
* @return (et.Items, error)
**/
func (s *Database) command(command *Command) (et.Items, error) {
	if err := command.validate(); err != nil {
		return et.Items{}, err
	}

	result, err := s.driver.Command(command)
	if err != nil {
		return et.Items{}, err
	}

	if command.isDebug {
		console.Debugf("command:%s", command.toJson().ToString())
	}

	return result, nil
}
