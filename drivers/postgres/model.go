package postgres

import (
	"fmt"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

func tableName(model *jdb.Model) string {
	return fmt.Sprintf(`%s.%s`, model.Schema, model.Name)
}

/**
* existTable
* @param schema string
* @param name string
* @return bool, error
**/
func (s *Postgres) existTable(schema, name string) (bool, error) {
	items, err := jdb.Query(s.jdb, `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.tables
		WHERE UPPER(table_schema) = UPPER($1)
		AND UPPER(table_name) = UPPER($2));`, schema, name)
	if err != nil {
		return false, err
	}

	if items.Count == 0 {
		return false, nil
	}

	return items.Bool(0, "exists"), nil
}

/**
* LoadModel
* @param model *jdb.Model
* @return error
**/
func (s *Postgres) LoadModel(model *jdb.Model) error {
	err := s.loadSchema(model.Schema)
	if err != nil {
		return err
	}

	exist, err := s.existTable(model.Schema, model.Name)
	if err != nil {
		return err
	}

	if !exist {
		sql := s.ddlTable(model)
		sqlIndex := s.ddlTableIndex(model)
		sql = strs.Append(sql, sqlIndex, "\n")
		if model.IsDebug {
			console.Debug(sql)
		}

		err = jdb.Definition(s.jdb, sql)
		if err != nil {
			return err
		}

		console.Logf("Model", "Create %s", tableName(model))

		return nil
	}

	if model.UseCore {
		return nil
	}

	sql := `
	SELECT
	a.attname AS column_name, 
	t.typname AS data_type,
	CASE 
		WHEN a.attlen > 0 THEN a.attlen
		WHEN a.attlen = -1 AND a.atttypmod > 0 THEN a.atttypmod - 4
		ELSE NULL
	END AS size
	FROM pg_catalog.pg_attribute a
	JOIN pg_catalog.pg_class c ON a.attrelid = c.oid
	JOIN pg_catalog.pg_namespace n ON c.relnamespace = n.oid
	JOIN pg_catalog.pg_type t ON a.atttypid = t.oid
	WHERE n.nspname = $1
	AND c.relname = $2
	AND a.attnum > 0
	AND NOT a.attisdropped;`

	items, err := jdb.Query(s.jdb, sql, model.Schema, model.Name)
	if err != nil {
		return err
	}

	for _, item := range items.Result {
		name := item.Str("column_name")
		dataType := item.Str("data_type")
		size := item.Int("size")
		typeData := s.strToTypeData(dataType, size)
		model.DefineColumn(name, typeData)
	}

	return nil
}

/**
* DropModel
* @param model *jdb.Model
* @return error
**/
func (s *Postgres) DropModel(model *jdb.Model) error {
	sql := s.ddlTableDrop(tableName(model))
	if model.IsDebug {
		console.Debug(sql)
	}

	return jdb.Definition(s.jdb, sql)
}

/**
* EmptyModel
* @param model *jdb.Model
* @return error
**/
func (s *Postgres) EmptyModel(model *jdb.Model) error {
	sql := s.ddlTableEmpty(tableName(model))
	if model.IsDebug {
		console.Debug(sql)
	}

	return jdb.Definition(s.jdb, sql)
}

/**
* MutateModel
* @param model *jdb.Model
* @return error
**/
func (s *Postgres) MutateModel(model *jdb.Model) error {
	backupTable := fmt.Sprintf(`%s_backup`, tableName(model))
	sql := "\n"
	sql = strs.Append(sql, s.ddlTableRename(tableName(model), backupTable), "\n")
	sql = strs.Append(sql, s.ddlTable(model), "\n")
	sql = strs.Append(sql, s.ddlTableInsertTo(model, backupTable), "\n\n")
	sql = strs.Append(sql, s.ddlTableIndex(model), "\n\n")
	sql = strs.Append(sql, s.ddlTableDrop(backupTable), "\n\n")
	if model.IsDebug {
		console.Debug(sql)
	}

	return jdb.Definition(s.jdb, sql)
}
