package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

/**
* Insert
* @param data []et.Json
* @return *Command
**/
func (s *Model) Insert(data et.Json) *Command {
	return NewCommand(s, []et.Json{data}, Insert)
}

/**
* Update
* @param data []et.Json
* @return *Command
**/
func (s *Model) Update(data et.Json) *Command {
	return NewCommand(s, []et.Json{data}, Update)
}

/**
* Delete
* @return *Command
**/
func (s *Model) Delete() *Command {
	return NewCommand(s, []et.Json{}, Delete)
}

/**
* Bulk
* @param data []et.Json
* @return *Command
**/
func (s *Model) Bulk(data []et.Json) *Command {
	return NewCommand(s, data, Bulk)
}

/**
* Undo
* @param index int64
* @return *Command
**/
func (s *Model) Undo(key string, index int64) *Command {
	result := NewCommand(s, []et.Json{}, Undo)
	result.Undo.Set("key", key)
	result.Undo.Set("index", index)

	return result
}

/**
* Command
* @param params et.Json
* @return interface{}, error
**/
func (s *Model) setCommand(params et.Json) (et.Items, error) {
	command := params.Str("command")
	var conm *Command
	switch command {
	case "insert":
		data := params.Json("data")
		conm = s.Insert(data)
	case "update":
		data := params.Json("data")
		conm = s.Update(data)
	case "delete":
		conm = s.Delete()
	case "bulk":
		data := params.ArrayJson("data")
		conm = s.Bulk(data)
	case "undo":
		key := params.Str("key")
		index := params.Int64("index")
		conm = s.Undo(key, index)
	}
	if conm == nil {
		return et.Items{}, mistake.New("command not found")
	}

	debug := params.ValBool(false, "debug")
	if debug {
		conm.Debug()
	}
	where := params.Json("where")
	conm.setWheres(where, conm.getField)

	return conm.Exec()
}
