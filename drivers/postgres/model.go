package postgres

import (
	"encoding/json"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) LoadModel(model *jdb.Model) error {
	current, err := s.getModel(model.Table)
	if err != nil {
		return err
	}

	var action string
	var sql string
	version := current.Int("version")
	exists, err := s.tableExists(model.Schema.Name, model.Table)
	if err != nil {
		return err
	}

	if exists {
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

			sql = s.ddlMutate(&old, model)
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

	if model.Show {
		console.Debug(sql)
	}

	for _, detail := range model.Details {
		err = s.LoadModel(detail.Model)
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
* LoadByTable
* @param model *jdb.Model
* @return error
**/
func (s *Postgres) LoadByTable(model *jdb.Model) error {
	sql := `
	SELECT column_name, data_type, character_maximum_length, is_nullable, column_default
	FROM information_schema.columns
	WHERE table_schema = $1 
  AND table_name = $2;`

	items, err := s.Query(sql, model.Schema.Name, model.Name)
	if err != nil {
		return nil
	}

	for _, item := range items.Result {
		name := item.Str("column_name")
		dataType := item.Str("data_type")
		typeData := s.strToTypeData(dataType)
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
	for _, detail := range model.Details {
		err := s.DropModel(detail.Model)
		if err != nil {
			return err
		}
	}

	return s.deleteModel(model.Table)
}
