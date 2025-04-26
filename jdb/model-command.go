package jdb

import (
	"github.com/cgalvisleon/et/et"
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
* Sync
* @return *Command
**/
func (s *Model) Sync() *Command {
	return NewCommand(s, []et.Json{}, Sync)
}
