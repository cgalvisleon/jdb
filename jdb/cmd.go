package jdb

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	"github.com/cgalvisleon/et/vm"
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
	Command           string           `json:"command"`
	From              *Model           `json:"from"`
	Data              []et.Json        `json:"data"`
	Result            et.Items         `json:"result"`
	UseAtribs         bool             `json:"use_atribs"`
	SQL               string           `json:"sql"`
	db                *Database        `json:"-"`
	tx                *Tx              `json:"-"`
	vm                *vm.Vm           `json:"-"`
	isDebug           bool             `json:"-"`
	beforeInserts     []string         `json:"-"`
	beforeUpdates     []string         `json:"-"`
	beforeDeletes     []string         `json:"-"`
	afterInserts      []string         `json:"-"`
	afterUpdates      []string         `json:"-"`
	afterDeletes      []string         `json:"-"`
	eventBeforeInsert []DataFunctionTx `json:"-"`
	eventBeforeUpdate []DataFunctionTx `json:"-"`
	eventBeforeDelete []DataFunctionTx `json:"-"`
	eventAfterInsert  []DataFunctionTx `json:"-"`
	eventAfterUpdate  []DataFunctionTx `json:"-"`
	eventAfterDelete  []DataFunctionTx `json:"-"`
}

/**
* newCommand
* @param model *Model, cmd string, data et.Json
* @return *Cmd
**/
func newCommand(model *Model, cmd string, data []et.Json) *Cmd {
	result := &Cmd{
		where:             newWhere(),
		Command:           cmd,
		From:              model,
		Data:              data,
		Result:            et.Items{},
		db:                model.db,
		vm:                vm.NewVm(),
		beforeInserts:     model.BeforeInserts,
		beforeUpdates:     model.BeforeUpdates,
		beforeDeletes:     model.BeforeDeletes,
		afterInserts:      model.AfterInserts,
		afterUpdates:      model.AfterUpdates,
		afterDeletes:      model.AfterDeletes,
		eventBeforeInsert: model.eventBeforeInsert,
		eventBeforeUpdate: model.eventBeforeUpdate,
		eventBeforeDelete: model.eventBeforeDelete,
		eventAfterInsert:  model.eventAfterInsert,
		eventAfterUpdate:  model.eventAfterUpdate,
		eventAfterDelete:  model.eventAfterDelete,
	}
	result.UseAtribs = model.SourceField != "" && !model.IsLocked

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
* EventBeforeInsert
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) EventBeforeInsert(fn DataFunctionTx) *Cmd {
	s.eventBeforeInsert = append(s.eventBeforeInsert, fn)
	return s
}

/**
* EventBeforeUpdate
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) EventBeforeUpdate(fn DataFunctionTx) *Cmd {
	s.eventBeforeUpdate = append(s.eventBeforeUpdate, fn)
	return s
}

/**
* EventBeforeDelete
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) EventBeforeDelete(fn DataFunctionTx) *Cmd {
	s.eventBeforeDelete = append(s.eventBeforeDelete, fn)
	return s
}

/**
* EventAfterInsert
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) EventAfterInsert(fn DataFunctionTx) *Cmd {
	s.eventAfterInsert = append(s.eventAfterInsert, fn)
	return s
}

/**
* EventAfterUpdate
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) EventAfterUpdate(fn DataFunctionTx) *Cmd {
	s.eventAfterUpdate = append(s.eventAfterUpdate, fn)
	return s
}

/**
* EventAfterDelete
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) EventAfterDelete(fn DataFunctionTx) *Cmd {
	s.eventAfterDelete = append(s.eventAfterDelete, fn)
	return s
}

/**
* EventBeforeInsertOrUpdate
* @param fn DataFunctionTx
* @return *Cmd
**/
func (s *Cmd) EventBeforeInsertOrUpdate(fn DataFunctionTx) *Cmd {
	s.eventBeforeInsert = append(s.eventBeforeInsert, fn)
	s.eventBeforeUpdate = append(s.eventBeforeUpdate, fn)
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
