package jdb

import "github.com/cgalvisleon/et/et"

/**
* Insert
* @param data et.Json
* @return *Cmd
**/
func (s *Model) Insert(data et.Json) *Cmd {
	return newCommand(s, TypeInsert, []et.Json{data})
}

/**
* Update
* @param data et.Json
* @return *Cmd
**/
func (s *Model) Update(data et.Json) *Cmd {
	return newCommand(s, TypeUpdate, []et.Json{data})
}

/**
* Delete
* @param data et.Json
* @return *Cmd
**/
func (s *Model) Delete(data et.Json) *Cmd {
	return newCommand(s, TypeDelete, []et.Json{data})
}

/**
* Upsert
* @param data et.Json
* @return *Cmd
**/
func (s *Model) Upsert(data et.Json) *Cmd {
	return newCommand(s, TypeUpsert, []et.Json{data})
}
