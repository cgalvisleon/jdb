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
func (s *Model) delete() *Command {
	return NewCommand(s, []et.Json{}, Delete)
}

/**
* Delete
* @param val string
* @return *Command
**/
func (s *Model) Delete(val string) *Command {
	result := s.delete()
	result.Where(val)

	return result
}

/**
* Upsert
* @param data []et.Json
* @return *Command
**/
func (s *Model) Upsert(data et.Json) *Command {
	return NewCommand(s, []et.Json{data}, Upsert)
}

/**
* Delsert
* @param data []et.Json
* @return *Command
**/
func (s *Model) Delsert(data et.Json) *Command {
	return NewCommand(s, []et.Json{data}, Delsert)
}

/**
* Bulk
* @param data []et.Json
* @return *Command
**/
func (s *Model) Bulk(data []et.Json) *Command {
	return NewCommand(s, data, Insert)
}

/**
* Sync
* @param data et.Json
* @param sysId string
* @return *Command
**/
func (s *Model) Sync(data et.Json, sysId string) *Command {
	result := NewCommand(s, []et.Json{data}, Sync)
	result.Where(cf.SystemId).Eq(sysId)
	result.isSync = true

	return result
}
