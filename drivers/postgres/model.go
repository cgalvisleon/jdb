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

	oldVersion := current.Int("version")
	var action = "load"
	var sql string
	if oldVersion == 0 {
		sql = s.ddlTable(model)
	} else if oldVersion != model.Version {
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
		action = "mutate"
	} else {
		sql = s.ddlIndexFunction(model)
	}

	console.Debug("LoadModel:", sql)

	serialized, err := model.Serialized()
	if err != nil {
		return err
	}

	err = s.Exec(sql)
	if err != nil {
		return err
	}

	go s.upsertDDL(strs.Format(`load_model_%s`, model.Table), sql)
	go s.upsertModel(model.Table, model.Version, serialized)

	console.Logf(jdb.Postgres, `Model %s %s`, model.Name, action)

	return nil
}
