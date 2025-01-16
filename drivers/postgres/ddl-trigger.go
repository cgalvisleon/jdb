package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func ddlRecordTriggers(model *jdb.Model) string {
	return defineRecordTrigger(model.Table)
}

func ddlRecycligTriggers(model *jdb.Model) string {
	return defineRecyclingTrigger(model.Table)
}

func ddlSeriesTriggers(model *jdb.Model) string {
	return defineSeriesTrigger(model.Table)
}

func (s *Postgres) ddlTriggers(model *jdb.Model) string {
	var result string
	if !model.Db.UseCore {
		return result
	}

	if model.SystemKeyField != nil {
		result = strs.Append(result, ddlRecordTriggers(model), "\n\n")
	}
	if model.StateField != nil {
		result = strs.Append(result, ddlRecycligTriggers(model), "\n\n")
	}
	if model.IndexField != nil {
		result = strs.Append(result, ddlSeriesTriggers(model), "\n\n")
	}

	return result
}
