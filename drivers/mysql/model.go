package mysql

import (
	"fmt"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

func tableName(model *jdb.Model) string {
	return fmt.Sprintf(`%s_%s`, model.Schema, model.Name)
}

/**
* existTable
* @param db, name string
* @return bool, error
**/
func (s *Mysql) existTable(db, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.tables
		WHERE UPPER(table_schema) = UPPER($1)
		AND UPPER(table_name) = UPPER($2));`

	items, err := jdb.Query(s.jdb, sql, db, name)
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
func (s *Mysql) LoadModel(model *jdb.Model) error {
	table := tableName(model)
	exist, err := s.existTable(model.Db.Name, table)
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

		err = jdb.Ddl(s.jdb, sql)
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
	COLUMN_NAME,
  DATA_TYPE,
  CHARACTER_MAXIMUM_LENGTH AS max_length,
  NUMERIC_PRECISION,
  NUMERIC_SCALE
	FROM INFORMATION_SCHEMA.COLUMNS
	WHERE TABLE_SCHEMA = $1
	AND TABLE_NAME = $2;`

	items, err := jdb.Query(s.jdb, sql, s.jdb.Name, table)
	if err != nil {
		return err
	}

	for _, item := range items.Result {
		name := item.Str("column_name")
		dataType := item.Str("data_type")
		size := item.Int("max_length")
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
func (s *Mysql) DropModel(model *jdb.Model) error {
	sql := s.ddlTableDrop(tableName(model))
	if model.IsDebug {
		console.Debug(sql)
	}

	_, err := jdb.Query(s.jdb, sql)
	if err != nil {
		return err
	}

	return nil
}

/**
* EmptyModel
* @param model *jdb.Model
* @return error
**/
func (s *Mysql) EmptyModel(model *jdb.Model) error {
	sql := s.ddlTableEmpty(tableName(model))
	if model.IsDebug {
		console.Debug(sql)
	}

	_, err := jdb.Query(s.jdb, sql)
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
func (s *Mysql) MutateModel(model *jdb.Model) error {
	backupTable := strs.Format(`%s_backup`, tableName(model))
	sql := "\n"
	sql = strs.Append(sql, s.ddlTableRename(tableName(model), backupTable), "\n")
	sql = strs.Append(sql, s.ddlTable(model), "\n")
	sql = strs.Append(sql, s.ddlTableInsertTo(model, backupTable), "\n\n")
	sql = strs.Append(sql, s.ddlTableIndex(model), "\n\n")
	sql = strs.Append(sql, s.ddlTableDrop(backupTable), "\n\n")
	if model.IsDebug {
		console.Debug(sql)
	}

	_, err := jdb.Query(s.jdb, sql)
	if err != nil {
		return err
	}

	return nil
}
