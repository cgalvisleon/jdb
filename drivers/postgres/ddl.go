package postgres

import (
	"fmt"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	"github.com/cgalvisleon/jdb/jdb"
)

/**
* buildModel
* @param definition et.Json
* @return (string, error)
**/
func (s *Postgres) buildModel(definition et.Json) (string, error) {
	sql, err := s.buildSchema(definition)
	if err != nil {
		return "", err
	}

	def, err := s.buildTable(definition)
	if err != nil {
		return "", err
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n\t")
	}

	def, err = s.buildPrimaryKeys(definition)
	if err != nil {
		return "", err
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n\t")
	}

	def, err = s.buildForeignKeys(definition)
	if err != nil {
		return "", err
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n\t")
	}

	def, err = s.buildIndices(definition)
	if err != nil {
		return "", err
	}

	if def != "" {
		sql = strs.Append(sql, def, "\n\t")
	}

	return sql, nil
}

/**
* buildSchema
* @param definition et.Json
* @return (string, error)
**/
func (s *Postgres) buildSchema(definition et.Json) (string, error) {
	schema := definition.String("schema")
	if !utility.ValidStr(schema, 0, []string{}) {
		return "", fmt.Errorf(jdb.MSG_SCHEMA_REQUIRED)
	}

	return fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", schema), nil
}

/**
* buildTable
* @param definition et.Json
* @return (string, error)
**/
func (s *Postgres) buildTable(definition et.Json) (string, error) {
	getType := func(tp string) string {
		switch tp {
		case "int":
			return "BIGINT"
		case "float":
			return "DOUBLE PRECISION"
		case "key":
			return "VARCHAR(80)"
		case "text":
			return "TEXT"
		case "datetime":
			return "TIMESTAMP"
		case "boolean":
			return "BOOLEAN"
		case "json":
			return "JSONB"
		case "index":
			return "BIGINT"
		case "bytes":
			return "BYTEA"
		case "geometry":
			return "JSONB"
		default:
			return tp
		}
	}

	columns := definition.ArrayJson("columns")
	columnsDef := ""
	for _, v := range columns {
		tp := v.String("type")
		if !jdb.TypeColumn[tp] {
			continue
		}
		tp = getType(tp)
		def := fmt.Sprintf("\n\t%s %s", v.String("name"), tp)
		columnsDef = strs.Append(columnsDef, def, ",")
	}

	table := definition.String("table")
	result := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", table, columnsDef)

	return result, nil
}

/**
* buildPrimaryKeys
* @param definition et.Json
* @return (string, error)
**/
func (s *Postgres) buildPrimaryKeys(definition et.Json) (string, error) {
	primaryKeys := definition.ArrayStr("primary_keys")
	if len(primaryKeys) == 0 {
		return "", nil
	}

	columns := ""
	for _, v := range primaryKeys {
		columns = strs.Append(columns, v, ", ")
	}

	table := definition.String("table")
	name := definition.String("name")
	result := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s_pk PRIMARY KEY (%s);", table, name, columns)

	return result, nil
}

/**
* buildForeignKeys
* @param definition et.Json
* @return (string, error)
**/
func (s *Postgres) buildForeignKeys(definition et.Json) (string, error) {
	foreignKeys := definition.ArrayJson("foreign_keys")
	if len(foreignKeys) == 0 {
		return "", nil
	}

	console.Debug(foreignKeys)

	// table := definition.String("table")
	// name := definition.String("name")
	// for _, v := range foreignKeys {
	// 	references := v.Json("references")
	// 	columns := references.ArrayJson("columns")
	// 	onDelete := references.String("on_delete")
	// 	onUpdate := references.String("on_update")

	// 	columns = strs.Append(columns, v.String("name"), ", ")
	// }

	// result := fmt.Sprintf("ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s_fk FOREIGN KEY (%s) REFERENCES %s(%s);", table, name, columns)

	return "", nil
}

/**
* buildIndices
* @param definition et.Json
* @return (string, error)
**/
func (s *Postgres) buildIndices(definition et.Json) (string, error) {
	indices := definition.ArrayStr("indices")
	if len(indices) == 0 {
		return "", nil
	}

	table := definition.String("table")
	name := definition.String("name")
	result := ""
	for _, v := range indices {
		def := fmt.Sprintf("%s_%s_idx", name, v)
		def = fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s (%s);", def, table, v)
		result = strs.Append(result, def, "\n\t")
	}

	return result, nil
}
