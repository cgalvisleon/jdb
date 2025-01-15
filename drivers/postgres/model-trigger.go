package postgres

import (
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
