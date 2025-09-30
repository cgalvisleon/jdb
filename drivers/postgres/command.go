package postgres

import (
	"fmt"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

/**
* buildCommand
* @param cmd *jdb.Cmd
* @return (string, error)
**/
func (s *Postgres) buildCommand(cmd *jdb.Cmd) (string, error) {
	console.Debug("command:", cmd.ToJson().ToString())
	command := cmd.Command
	if !jdb.Commands[command] {
		return "", fmt.Errorf("command %s no soportado", command)
	}

	switch command {
	case jdb.CmdInsert:
		return s.buildInsert(cmd)
	case jdb.CmdUpdate:
		return s.buildUpdate(cmd)
	case jdb.CmdDelete:
		return s.buildDelete(cmd)
	}

	return "", nil
}

/**
* buildInsert
* @param cmd *jdb.Cmd
* @return (string, error)
**/
func (s *Postgres) buildInsert(cmd *jdb.Cmd) (string, error) {
	table := cmd.From.Table
	data := cmd.Data[0]
	into := ""
	values := ""
	atribs := et.Json{}
	returning := fmt.Sprintf(`to_jsonb(%s.*) AS result`, table)
	for k, v := range data {
		val := fmt.Sprintf(`%v`, jdb.Quote(v))
		_, ok := cmd.From.GetColumn(k)
		if ok {
			into = strs.Append(into, k, ", ")
			values = strs.Append(values, val, ", ")
			continue
		}

		if cmd.UseAtribs {
			atribs[k] = val
		}
	}

	if cmd.UseAtribs {
		into = strs.Append(into, cmd.From.SourceField, ", ")
		values = strs.Append(values, fmt.Sprintf(`'%v'::jsonb`, atribs.ToString()), ", ")
		returning = fmt.Sprintf("to_jsonb(A) - '%s'", cmd.From.SourceField)
	}

	sql := fmt.Sprintf("INSERT INTO %s(%s)\nVALUES(%s)\nRETURNING %s;", table, into, values, returning)
	return sql, nil
}

/**
* buildUpdate
* @param cmd *jdb.Cmd
* @return (string, error)
**/
func (s *Postgres) buildUpdate(cmd *jdb.Cmd) (string, error) {
	return "", nil
}

/**
* buildDelete
* @param cmd *jdb.Cmd
* @return (string, error)
**/
func (s *Postgres) buildDelete(cmd *jdb.Cmd) (string, error) {
	return "", nil
}
