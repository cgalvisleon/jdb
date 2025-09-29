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

type Cmd struct {
	*where
	Command       string           `json:"command"`
	Data          []et.Json        `json:"items"`
	Columns       et.Json          `json:"columns"`
	Atribs        et.Json          `json:"atrib"`
	Keys          et.Json          `json:"keys"`
	Before        et.Json          `json:"before"`
	After         et.Json          `json:"after"`
	Results       []et.Json        `json:"results"`
	SQL           string           `json:"sql"`
	db            *Database        `json:"-"`
	tx            *Tx              `json:"-"`
	from          *Model           `json:"-"`
	result        []et.Json        `json:"-"`
	isDebug       bool             `json:"-"`
	useAtribs     bool             `json:"-"`
	beforeInsert  []DataFunctionTx `json:"-"`
	beforeUpdate  []DataFunctionTx `json:"-"`
	beforeDelete  []DataFunctionTx `json:"-"`
	afterInsert   []DataFunctionTx `json:"-"`
	afterUpdate   []DataFunctionTx `json:"-"`
	afterDelete   []DataFunctionTx `json:"-"`
	beforeInserts []string         `json:"-"`
	beforeUpdates []string         `json:"-"`
	beforeDeletes []string         `json:"-"`
	afterInserts  []string         `json:"-"`
	afterUpdates  []string         `json:"-"`
	afterDeletes  []string         `json:"-"`
}

/**
* newCommand
* @param model *Model, cmd string, data []et.Json
* @return *Cmd
**/
func newCommand(model *Model, cmd string, data []et.Json) *Cmd {
	result := &Cmd{
		where:         newWhere(),
		Command:       cmd,
		Data:          data,
		Columns:       et.Json{},
		Atribs:        et.Json{},
		Keys:          et.Json{},
		Before:        et.Json{},
		After:         et.Json{},
		Results:       []et.Json{},
		db:            model.db,
		from:          model,
		beforeInserts: model.BeforeInserts,
		beforeUpdates: model.BeforeUpdates,
		beforeDeletes: model.BeforeDeletes,
		afterInserts:  model.AfterInserts,
		afterUpdates:  model.AfterUpdates,
		afterDeletes:  model.AfterDeletes,
		beforeInsert:  model.beforeInsert,
		beforeUpdate:  model.beforeUpdate,
		beforeDelete:  model.beforeDelete,
		afterInsert:   model.afterInsert,
		afterUpdate:   model.afterUpdate,
		afterDelete:   model.afterDelete,
	}
	result.useAtribs = model.SourceField != "" && !model.IsLocked

	return result
}

/**
* toJson
* @return et.Json
**/
func (s *Cmd) toJson() et.Json {
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
* @return *Cmd
**/
func (s *Cmd) Debug() *Cmd {
	s.isDebug = true
	return s
}

/**
* setTx
* @param tx *Tx
* @return *Ql
**/
func (s *Cmd) setTx(tx *Tx) *Cmd {
	s.tx = tx
	return s
}

/**
* command
* @param cmd string, param et.Json
* @return (*Cmd, error)
**/
func command(cmd string, param et.Json) (*Cmd, error) {
	database := param.String("database")
	if !utility.ValidStr(database, 0, []string{}) {
		return nil, fmt.Errorf(MSG_DATABASE_REQUIRED)
	}

	db, err := GetDatabase(database)
	if err != nil {
		return nil, err
	}

	name := param.String("name")
	if !utility.ValidStr(name, 0, []string{}) {
		return nil, fmt.Errorf(MSG_NAME_REQUIRED)
	}

	model, err := db.getModel(name)
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
* BeforeInsert
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) BeforeInsert(fn DataFunctionTx) *Cmd {
	s.beforeInsert = append(s.beforeInsert, fn)
	return s
}

/**
* BeforeUpdate
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) BeforeUpdate(fn DataFunctionTx) *Cmd {
	s.beforeUpdate = append(s.beforeUpdate, fn)
	return s
}

/**
* BeforeDelete
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) BeforeDelete(fn DataFunctionTx) *Cmd {
	s.beforeDelete = append(s.beforeDelete, fn)
	return s
}

/**
* AfterInsert
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) AfterInsert(fn DataFunctionTx) *Cmd {
	s.afterInsert = append(s.afterInsert, fn)
	return s
}

/**
* AfterUpdate
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) AfterUpdate(fn DataFunctionTx) *Cmd {
	s.afterUpdate = append(s.afterUpdate, fn)
	return s
}

/**
* AfterDelete
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) AfterDelete(fn DataFunctionTx) *Cmd {
	s.afterDelete = append(s.afterDelete, fn)
	return s
}

/**
* BeforeInsertOrUpdate
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) BeforeInsertOrUpdate(fn DataFunctionTx) *Cmd {
	s.beforeInsert = append(s.beforeInsert, fn)
	s.beforeUpdate = append(s.beforeUpdate, fn)
	return s
}

/**
* Insert
* @param param et.Json
* @return (*Cmd, error)
**/
func Insert(param et.Json) (*Cmd, error) {
	return command(TypeInsert, param)
}

/**
* Update
* @param param et.Json
* @return (*Cmd, error)
**/
func Update(param et.Json) (*Cmd, error) {
	return command(TypeUpdate, param)
}

/**
* Delete
* @param param et.Json
* @return (*Cmd, error)
**/
func Delete(param et.Json) (*Cmd, error) {
	return command(TypeDelete, param)
}

/**
* Upsert
* @param param et.Json
* @return (*Cmd, error)
**/
func Upsert(param et.Json) (*Cmd, error) {
	return command(TypeUpsert, param)
}
