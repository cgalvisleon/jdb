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
	result := NewCommand(s, []et.Json{data}, Insert)
	result.Show = s.Show

	return result
}

/**
* Update
* @param data []et.Json
* @return *Command
**/
func (s *Model) Update(data et.Json) *Command {
	result := NewCommand(s, []et.Json{data}, Update)
	result.Show = s.Show

	return result
}

/**
* Delete
* @return *Command
**/
func (s *Model) Delete() *Command {
	result := NewCommand(s, []et.Json{}, Delete)
	result.Show = s.Show

	return result
}

/**
* Bulk
* @param data []et.Json
* @return *Command
**/
func (s *Model) Bulk(data []et.Json) *Command {
	result := NewCommand(s, data, Bulk)
	result.Show = s.Show

	return result
}

/**
* Undo
* @param key string
* @param index int64
* @return *Command
**/
func (s *Model) Undo(key string, index int64) *Command {
	result := NewCommand(s, []et.Json{}, Undo)
	result.Show = s.Show
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
