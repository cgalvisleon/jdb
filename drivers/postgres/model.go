package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) CreateModel(model *jdb.Model) error {
	sql := s.ddlTable(model)

	console.Debug("CreateModel:", sql)
	err := s.Exec(sql)
	if err != nil {
		return mistake.Newf(jdb.MSG_QUERY_FAILED, err.Error())
	}

	go s.upsertDDL(strs.Format(`create_model_%s`, model.Table), sql)

	console.Logf(jdb.Postgres, `Model %s created`, model.Name)

	return nil
}

func (s *Postgres) MutateModel(model *jdb.Model) error {
	sql := s.ddlTable(model)
	err := s.Exec(sql)
	if err != nil {
		return mistake.Newf(jdb.MSG_QUERY_FAILED, err.Error())
	}

	go s.upsertDDL(strs.Format(`mutate_model_%s`, model.Table), sql)

	console.Logf(jdb.Postgres, `Model %s mutated`, model.Name)

	return nil
}
