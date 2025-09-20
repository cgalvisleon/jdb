package jdb

import (
	"encoding/json"
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

const (
	TypeInsert = "insert"
	TypeUpdate = "update"
	TypeDelete = "delete"
	TypeUpsert = "upsert"
)

var (
	Commands map[string]bool = map[string]bool{
		TypeInsert: true,
		TypeUpdate: true,
		TypeDelete: true,
		TypeUpsert: true,
	}
)

type Command struct {
	Command string    `json:"command"`
	Data    []et.Json `json:"data"`
	Before  et.Json   `json:"before"`
	After   et.Json   `json:"after"`
	SQL     string    `json:"sql"`
	db      *Database `json:"-"`
	tx      *Tx       `json:"-"`
	model   *Model    `json:"-"`
	isDebug bool      `json:"-"`
}

/**
* newCommand
* @param model *Model, cmd string, data []et.Json
* @return *Command
**/
func newCommand(model *Model, cmd string, data []et.Json) *Command {
	return &Command{
		Command: cmd,
		Data:    data,
		Before:  et.Json{},
		After:   et.Json{},
		db:      model.db,
		model:   model,
	}
}

/**
* toJson
* @return et.Json
**/
func (s *Command) toJson() et.Json {
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
* Debug
* @return *Command
**/
func (s *Command) Debug() *Command {
	s.isDebug = true
	return s
}

/**
* Exec
* @param tx *Tx
* @return (et.Items, error)
**/
func (s *Command) Exec(tx *Tx) (et.Items, error) {
	s.tx = tx

	if err := s.validate(); err != nil {
		return et.Items{}, err
	}

	return et.Items{}, nil
}

/**
* command
* @param cmd string, param et.Json
* @return (*Command, error)
**/
func command(cmd string, param et.Json) (*Command, error) {
	database := param.String("database")
	if !utility.ValidStr(database, 0, []string{}) {
		return nil, fmt.Errorf(MSG_DATABASE_REQUIRED)
	}

	db, err := GetDatabase(database)
	if err != nil {
		return nil, err
	}

	eschema := param.String("schema")
	if !utility.ValidStr(eschema, 0, []string{}) {
		return nil, fmt.Errorf(MSG_SCHEMA_REQUIRED)
	}

	name := param.String("name")
	if !utility.ValidStr(name, 0, []string{}) {
		return nil, fmt.Errorf(MSG_NAME_REQUIRED)
	}

	model, err := db.getOrCreateModel(eschema, name)
	if err != nil {
		return nil, err
	}

	data := param.ArrayJson("data")
	if len(data) == 0 {
		return nil, fmt.Errorf(MSG_DATA_REQUIRED)
	}

	return newCommand(model, cmd, data), nil
}

/**
* Insert
* @param param et.Json
* @return (*Command, error)
**/
func Insert(param et.Json) (*Command, error) {
	return command(TypeInsert, param)
}

/**
* Update
* @param param et.Json
* @return (*Command, error)
**/
func Update(param et.Json) (*Command, error) {
	return command(TypeUpdate, param)
}

/**
* Delete
* @param param et.Json
* @return (*Command, error)
**/
func Delete(param et.Json) (*Command, error) {
	return command(TypeDelete, param)
}

/**
* Upsert
* @param param et.Json
* @return (*Command, error)
**/
func Upsert(param et.Json) (*Command, error) {
	return command(TypeUpsert, param)
}
