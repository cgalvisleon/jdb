package postgres

import (
	"encoding/json"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

/**
* LoadTable
* @param model *jdb.Model
* @return error
**/
func (s *Postgres) LoadTable(model *jdb.Model) (bool, error) {
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
		return false, err
	}

	for _, item := range items.Result {
		name := item.Str("column_name")
		dataType := item.Str("data_type")
		size := item.Int("size")
		typeData := s.strToTypeData(dataType, size)
		model.DefineColumn(name, typeData)
	}

	return items.Ok, nil
}

/**
* LoadModel
* @param model *jdb.Model
* @return error
**/
func (s *Postgres) LoadModel(model *jdb.Model) error {
	current, err := s.getModel(model.Table)
	if err != nil {
		return err
	}

	var action string
	var sql string
	version := current.Int("version")
	if model.IsCreated {
		if version != model.Version {
			action = "mutate"
			bt, err := current.Byte("model")
			if err != nil {
				return err
			}
			var old jdb.Model
			err = json.Unmarshal(bt, &old)
			if err != nil {
				return err
			}

			sql = s.ddlMutate(&old, model, false)
		} else {
			action = "index"
			sql = s.ddlIndexFunction(model)
		}
	} else {
		action = "load"
		sql = s.ddlTable(model)
	}

	serialized, err := model.Serialized()
	if err != nil {
		return err
	}

	id := strs.Format(`load_model_%s`, model.Table)
	err = s.ExecDDL(id, sql)
	if err != nil {
		return err
	}

	if model.IsDebug {
		console.Debug(sql)
	}

	for _, detail := range model.Details {
		err = s.LoadModel(detail.With)
		if err != nil {
			model.Drop()
			return err
		}
	}

	go s.upsertModel(model.Table, model.Version, serialized)

	console.Logf(model.Db.Name, `Model %s %s`, model.Name, action)

	return nil
}

/**
* DropModel
* @param model *jdb.Model
* @return error
**/
func (s *Postgres) DropModel(model *jdb.Model) error {
	for _, detail := range model.Details {
		err := s.DropModel(detail.With)
		if err != nil {
			return err
		}
	}

	return s.deleteModel(model.Table)
}
