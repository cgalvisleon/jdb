package postgres

import (
	"encoding/json"

	"github.com/cgalvisleon/et/console"
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

	console.Debug("LoadModel:", sql)

	// id := strs.Format(`load_model_%s`, model.Table)
	// err = s.ExecDDL(id, sql)
	// if err != nil {
	// 	return err
	// }

	go s.upsertModel(model.Table, model.Version, serialized)

	console.Logf(model.Db.Name, `Model %s %s`, model.Name, action)

	return nil
}
