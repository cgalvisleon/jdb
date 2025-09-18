package jdb

import (
	"encoding/json"

	"github.com/cgalvisleon/et/et"
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
