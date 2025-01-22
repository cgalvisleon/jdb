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
* @param key string
* @param index int64
* @return *Command
**/
func (s *Model) Undo(key string, index int64) *Command {
	result := NewCommand(s, []et.Json{}, Undo)
	result.Undo = &UndoRecord{
		Key:   key,
		Index: index,
	}

	return result
}

/**
* Command
* @param params et.Json
* @return et.Items, error
**/
func (s *Model) Command(params et.Json) (et.Items, error) {
	command := params.Str("command")
	where := params.ArrayJson("where")
	debug := params.ValBool(false, "debug")
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
	}
	if conm == nil {
		return et.Items{}, mistake.New("command not found")
	}

	conm.setWhere(where)
	if debug {
		conm.Debug()
	}
	conm.Exec()
	return et.Items{
		Ok: true,
		Result: []et.Json{{
			"command": command,
			"from":    s.Table,
			"where":   conm.listWheres(),
		}},
	}, nil
}
