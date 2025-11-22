package jdb

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	"github.com/cgalvisleon/workflow/vm"
)

const (
	CmdInsert = "insert"
	CmdUpdate = "update"
	CmdDelete = "delete"
	CmdUpsert = "upsert"
)

var (
	Commands map[string]bool = map[string]bool{
		CmdInsert: true,
		CmdUpdate: true,
		CmdDelete: true,
		CmdUpsert: true,
	}
)

type Cmd struct {
	*where
	Command       string           `json:"command"`
	Data          []et.Json        `json:"data"`
	Result        et.Items         `json:"result"`
	SQL           string           `json:"sql"`
	db            *DB              `json:"-"`
	tx            *Tx              `json:"-"`
	vm            *vm.Vm           `json:"-"`
	IsDebug       bool             `json:"-"`
	BeforeInserts []string         `json:"-"`
	BeforeUpdates []string         `json:"-"`
	BeforeDeletes []string         `json:"-"`
	AfterInserts  []string         `json:"-"`
	AfterUpdates  []string         `json:"-"`
	AfterDeletes  []string         `json:"-"`
	beforeInserts []DataFunctionTx `json:"-"`
	beforeUpdates []DataFunctionTx `json:"-"`
	beforeDeletes []DataFunctionTx `json:"-"`
	afterInserts  []DataFunctionTx `json:"-"`
	afterUpdates  []DataFunctionTx `json:"-"`
	afterDeletes  []DataFunctionTx `json:"-"`
}

/**
* newCommand
* @param model *Model, cmd string, data et.Json
* @return *Cmd
**/
func newCommand(model *Model, cmd string, data []et.Json) *Cmd {
	result := &Cmd{
		where:         newWhere(model, ""),
		Command:       cmd,
		Data:          data,
		Result:        et.Items{},
		db:            model.db,
		vm:            vm.New(),
		BeforeInserts: model.BeforeInserts,
		BeforeUpdates: model.BeforeUpdates,
		BeforeDeletes: model.BeforeDeletes,
		AfterInserts:  model.AfterInserts,
		AfterUpdates:  model.AfterUpdates,
		AfterDeletes:  model.AfterDeletes,
		beforeInserts: model.beforeInserts,
		beforeUpdates: model.beforeUpdates,
		beforeDeletes: model.beforeDeletes,
		afterInserts:  model.afterInserts,
		afterUpdates:  model.afterUpdates,
		afterDeletes:  model.afterDeletes,
	}
	result.From = model

	return result
}

/**
* ToJson
* @return et.Json
**/
func (s *Cmd) ToJson() et.Json {
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
	s.IsDebug = true
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
* prepare
* @param data et.Json, as string
* @return []et.Json
**/
func (s *Cmd) getKeys(data et.Json, as string) []et.Json {
	result := []et.Json{}
	for k, v := range data {
		filed := strs.Append(as, k, ".")
		_, ok := s.From.GetColumn(k)
		if !ok {
			continue
		}

		if slices.Contains(s.From.PrimaryKeys, k) {
			result = append(result, et.Json{
				filed: et.Json{
					"eq": Quote(v),
				},
			})
		}
	}

	return result
}

/**
* BeforeInsert
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) BeforeInsert(fn DataFunctionTx) *Cmd {
	s.beforeInserts = append(s.beforeInserts, fn)
	return s
}

/**
* BeforeUpdate
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) BeforeUpdate(fn DataFunctionTx) *Cmd {
	s.beforeUpdates = append(s.beforeUpdates, fn)
	return s
}

/**
* BeforeDelete
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) BeforeDelete(fn DataFunctionTx) *Cmd {
	s.beforeDeletes = append(s.beforeDeletes, fn)
	return s
}

/**
* AfterInsert
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) AfterInsert(fn DataFunctionTx) *Cmd {
	s.afterInserts = append(s.afterInserts, fn)
	return s
}

/**
* AfterUpdate
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) AfterUpdate(fn DataFunctionTx) *Cmd {
	s.afterUpdates = append(s.afterUpdates, fn)
	return s
}

/**
* AfterDelete
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) AfterDelete(fn DataFunctionTx) *Cmd {
	s.afterDeletes = append(s.afterDeletes, fn)
	return s
}

/**
* BeforeInsertOrUpdate
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) BeforeInsertOrUpdate(fn DataFunctionTx) *Cmd {
	s.beforeInserts = append(s.beforeInserts, fn)
	s.beforeUpdates = append(s.beforeUpdates, fn)
	return s
}

/**
* AfterInsertOrUpdate
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) AfterInsertOrUpdate(fn DataFunctionTx) *Cmd {
	s.afterInserts = append(s.afterInserts, fn)
	s.afterUpdates = append(s.afterUpdates, fn)
	return s
}

/**
* Insert
* @param param et.Json
* @return (*Cmd, error)
**/
func Insert(param et.Json) (*Cmd, error) {
	return command(CmdInsert, param)
}

/**
* Update
* @param param et.Json
* @return (*Cmd, error)
**/
func Update(param et.Json) (*Cmd, error) {
	return command(CmdUpdate, param)
}

/**
* Delete
* @param param et.Json
* @return (*Cmd, error)
**/
func Delete(param et.Json) (*Cmd, error) {
	return command(CmdDelete, param)
}

/**
* Upsert
* @param param et.Json
* @return (*Cmd, error)
**/
func Upsert(param et.Json) (*Cmd, error) {
	return command(CmdUpsert, param)
}
