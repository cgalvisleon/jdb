package jdb

import "github.com/cgalvisleon/et/et"

/**
* Insert
* @param data et.Json
* @return *Cmd
**/
func (s *Model) Insert(data et.Json) *Cmd {
	return newCommand(s, CmdInsert, []et.Json{data})
}

/**
* Update
* @param data et.Json
* @return *Cmd
**/
func (s *Model) Update(data et.Json) *Cmd {
	return newCommand(s, CmdUpdate, []et.Json{data})
}

/**
* Delete
* @param data et.Json
* @return *Cmd
**/
func (s *Model) Delete() *Cmd {
	return newCommand(s, CmdDelete, []et.Json{})
}

/**
* Upsert
* @param data et.Json
* @return *Cmd
**/
func (s *Model) Upsert(data et.Json) *Cmd {
	return newCommand(s, CmdUpsert, []et.Json{data})
}
