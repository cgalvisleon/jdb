package sqlite

import (
	"fmt"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

/**
* existTable
* @param name string
* @return bool, error
**/
func (s *SqlLite) existTable(name string) (bool, error) {
	sql := `
	SELECT name
	FROM sqlite_master
	WHERE type='table'
	AND name=?;`

	items, err := jdb.Query(s.jdb, sql, name)
	if err != nil {
		return false, err
	}

	return items.Count > 0, nil
}

/**
* LoadModel
* @param model *jdb.Model
* @return error
**/
func (s *SqlLite) LoadModel(model *jdb.Model) error {
	model.Table = fmt.Sprintf(`%s_%s`, model.Schema, model.Name)
	exist, err := s.existTable(model.Table)
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

		console.Logf("Model", "Create %s", model.Table)

		return nil
	}

	if model.UseCore {
		return nil
	}

	sql := `
	SELECT 
    name AS column_name,
    type AS data_type,
    256 AS size
	FROM pragma_table_info(?);`

	items, err := jdb.Query(s.jdb, sql, model.Table)
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
func (s *SqlLite) DropModel(model *jdb.Model) error {
	sql := s.ddlTableDrop(model.Table)
	if model.IsDebug {
		console.Debug(sql)
	}

	err := jdb.Definition(s.jdb, sql)
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
func (s *SqlLite) EmptyModel(model *jdb.Model) error {
	sql := s.ddlTableEmpty(model.Table)
	if model.IsDebug {
		console.Debug(sql)
	}

	err := jdb.Definition(s.jdb, sql)
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
func (s *SqlLite) MutateModel(model *jdb.Model) error {
	backupTable := fmt.Sprintf(`%s_backup`, model.Table)
	sql := "\n"
	sql = strs.Append(sql, s.ddlTableRename(model.Table, backupTable), "\n")
	sql = strs.Append(sql, s.ddlTable(model), "\n")
	sql = strs.Append(sql, s.ddlTableInsertTo(model, backupTable), "\n\n")
	sql = strs.Append(sql, s.ddlTableIndex(model), "\n\n")
	sql = strs.Append(sql, s.ddlTableDrop(backupTable), "\n\n")
	if model.IsDebug {
		console.Debug(sql)
	}

	err := jdb.Definition(s.jdb, sql)
	if err != nil {
		return err
	}

	return nil
}
