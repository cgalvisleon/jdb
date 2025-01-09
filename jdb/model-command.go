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
* Command
* @param params et.Json
* @return et.Items, error
**/
func (s *Model) Command(params et.Json) (et.Items, error) {
	command := params.Str("command")
	where := params.ArrayJson([]et.Json{}, "where")
	returns := params.ArrayStr([]string{}, "returns")
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
		data := params.ArrayJson([]et.Json{}, "data")
		conm = s.Bulk(data)
	}
	if conm == nil {
		return et.Items{}, mistake.New("command not found")
	}

	conm.setWhere(where)
	conm.Return(returns...)
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
			"returns": conm.listReturns(),
		}},
	}, nil
}
