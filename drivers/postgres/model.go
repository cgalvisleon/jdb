package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

/**
* LoadModel
* @param model *jdb.Model
* @return error
**/
func (s *Postgres) LoadModel(model *jdb.Model) error {
	existTable, err := s.existTable(model.Schema.Name, model.Name)
	if err != nil {
		return err
	}

	if !existTable {
		sql := s.ddlTable(model)
		sqlIndex := s.ddlTableIndex(model)
		sql = strs.Append(sql, sqlIndex, "\n")
		if model.IsDebug {
			console.Debug(sql)
		}

		err = s.Exec(sql)
		if err != nil {
			return err
		}

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

	items, err := s.Query(sql, model.Schema.Name, model.Name)
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
	sql := s.ddlTableDrop(model.Table)
	if model.IsDebug {
		console.Debug(sql)
	}

	err := s.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

/**
* MutateModel
* @param model *jdb.Model
* @return error
**/
func (s *Postgres) MutateModel(model *jdb.Model) error {
	backupTable := strs.Format(`%s_backup`, model.Table)
	sql := "\n"
	sql = strs.Append(sql, s.ddlTableRename(model.Table, backupTable), "\n")
	sql = strs.Append(sql, s.ddlTable(model), "\n")
	sql = strs.Append(sql, s.ddlTableInsertTo(model, backupTable), "\n\n")
	sql = strs.Append(sql, s.ddlTableIndex(model), "\n\n")
	if model.IsDebug {
		console.Debug(sql)
	}
	err := s.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
