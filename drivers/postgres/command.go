package postgres

import (
	"fmt"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
)

/**
* buildCommand
* @param query et.Json
* @return (string, error)
**/
func (s *Postgres) buildCommand(query et.Json) (string, error) {
	console.Debug("command:", query.ToString())

	command := query.String("command")
	if !jdb.Commands[command] {
		return "", fmt.Errorf("command %s no soportado", command)
	}

	switch command {
	case jdb.CmdInsert:
		return s.buildInsert(query)
	case jdb.CmdUpdate:
		return s.buildUpdate(query)
	case jdb.CmdDelete:
		return s.buildDelete(query)
	}

	return "", nil
}

/**
* buildInsert
* @param query et.Json
* @return (string, error)
**/
func (s *Postgres) buildInsert(query et.Json) (string, error) {
	return "", nil
}

/**
* buildUpdate
* @param query et.Json
* @return (string, error)
**/
func (s *Postgres) buildUpdate(query et.Json) (string, error) {
	return "", nil
}

/**
* buildDelete
* @param query et.Json
* @return (string, error)
**/
func (s *Postgres) buildDelete(query et.Json) (string, error) {
	return "", nil
}
