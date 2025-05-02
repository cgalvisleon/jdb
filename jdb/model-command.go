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
* Upsert
* @param data []et.Json
* @return *Command
**/
func (s *Model) Upsert(data et.Json) *Command {
	return NewCommand(s, []et.Json{data}, Upsert)
}

/**
* Delete
* @param val string
* @return *Command
**/
func (s *Model) Delete(val string) *Command {
	result := NewCommand(s, []et.Json{}, Delete)
	result.Where(val)

	return result
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
* @param data et.Json, sysId string
* @return *Command
**/
func (s *Model) Undo(data et.Json, sysId string) *Command {
	result := NewCommand(s, []et.Json{}, Update)
	if s.History == nil {
		return result
	}

	if s.SystemKeyField == nil {
		return result
	}

	history := s.History.With
	if history == nil {
		return result
	}

	current, err := history.
		Where(SYSID).Eq(sysId).
		OrderByDesc(INDEX).
		One()
	if err != nil {
		return result
	}

	if !current.Ok {
		return result
	}

	historycal := current.Json(HISTORYCAL)
	if historycal.IsEmpty() {
		return result
	}

	result.Data = append(result.Data, historycal)
	result.Where(SYSID).Eq(sysId)
	result.isUndo = true

	return result
}

/**
* Sync
* @param data et.Json
* @param sysId string
* @return *Command
**/
func (s *Model) Sync(data et.Json, sysId string) *Command {
	result := NewCommand(s, []et.Json{data}, Sync)
	result.Where(SYSID).Eq(sysId)
	result.isSync = true

	return result
}
