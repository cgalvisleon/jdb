package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) typeColumn(tp jdb.TypeData) interface{} {
	switch tp {
	case jdb.TypeDataArray:
		return "JSONB"
	case jdb.TypeDataBool:
		return "BOOLEAN"
	case jdb.TypeDataInt:
		return "INTEGER"
	case jdb.TypeDataKey:
		return "VARCHAR(80)"
	case jdb.TypeDataMemo:
		return "TEXT"
	case jdb.TypeDataNumber:
		return "DECIMAL(18,2)"
	case jdb.TypeDataPrecision:
		return "DOUBLE PRECISION"
	case jdb.TypeDataObject:
		return "JSONB"
	case jdb.TypeDataSerie:
		return "BIGINT"
	case jdb.TypeDataShortText:
		return "VARCHAR(80)"
	case jdb.TypeDataText:
		return "VARCHAR(250)"
	case jdb.TypeDataTime:
		return "TIMESTAMP"
	default:
		return "VARCHAR(250)"
	}
}

func (s *Postgres) defaultValue(tp jdb.TypeData) interface{} {
	switch tp {
	case jdb.TypeDataArray:
		return "[]"
	case jdb.TypeDataBool:
		return "FALSE"
	case jdb.TypeDataInt:
		return 0
	case jdb.TypeDataKey:
		return "-1"
	case jdb.TypeDataMemo:
		return ""
	case jdb.TypeDataNumber:
		return 0.0
	case jdb.TypeDataObject:
		return "{}"
	case jdb.TypeDataSerie:
		return 0
	case jdb.TypeDataShortText:
		return ""
	case jdb.TypeDataText:
		return ""
	case jdb.TypeDataTime:
		return "NOW()"
	default:
		return ""
	}
}

func (s *Postgres) ddlIndex(model *jdb.Model) string {
	var result string
	for _, index := range model.Indices {
		def := ddlIndex(index.Column)

		result = strs.Append(result, def, "\n")
	}

	console.Debug("DDL Index:", result)
	return result
}

func (s *Postgres) ddlUniqueIndex(model *jdb.Model) string {
	var result string
	for _, index := range model.Uniques {
		def := ddlUniqueIndex(index.Column)

		result = strs.Append(result, def, "\n")
	}

	console.Debug("DDL Unique Index:", result)
	return result
}

func (s *Postgres) ddlTriggers(model *jdb.Model) string {
	var result string
	if !model.Db.UseCore {
		return result
	}

	if model.SystemKeyField != nil {
		result = strs.Append(result, ddlRecordTriggers(model), "\n")
	}
	if model.StateField != nil {
		result = strs.Append(result, ddlRecycligTriggers(model), "\n")
	}
	if model.IndexField != nil {
		result = strs.Append(result, ddlSeriesTriggers(model), "\n")
	}

	return result
}

func (s *Postgres) ddlTable(model *jdb.Model) string {
	var columnsDef string
	for _, column := range model.Columns {
		if column.TypeColumn == jdb.TpColumn {
			columnsDef += strs.Format("\t%s %s DEFAULT %v,\n", column.Name, s.typeColumn(column.TypeData), s.defaultValue(column.TypeData))
		}
	}
	columnsDef = strs.Append(columnsDef, ddlPrimaryKey(model), "\n")
	result := strs.Format("\nCREATE TABLE IF NOT EXISTS %s (\n%s);", model.Table, columnsDef)
	result = strs.Append(result, s.ddlIndex(model), "\n")
	result = strs.Append(result, s.ddlUniqueIndex(model), "\n")
	result = strs.Append(result, s.ddlTriggers(model), "\n")

	console.Debug("DDL Table:", result)
	return result
}

func (s *Postgres) CreateModel(model *jdb.Model) error {
	sql := s.ddlTable(model)

	console.Debug("CreateModel:", sql)
	/*
		err := s.Exec(sql)
		if err != nil {
			return mistake.Newf(jdb.MSG_QUERY_FAILED, err.Error())
		}

		go s.upsertDDL(strs.Format(`create_model_%s`, model.Table), sql)
	*/
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
