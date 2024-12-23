package postgres

import (
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
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
		return utility.Quote("[]")
	case jdb.TypeDataBool:
		return utility.Quote("FALSE")
	case jdb.TypeDataInt:
		return 0
	case jdb.TypeDataKey:
		return utility.Quote("-1")
	case jdb.TypeDataMemo:
		return utility.Quote("")
	case jdb.TypeDataNumber:
		return 0.0
	case jdb.TypeDataObject:
		return utility.Quote("{}")
	case jdb.TypeDataSerie:
		return 0
	case jdb.TypeDataShortText:
		return utility.Quote("")
	case jdb.TypeDataText:
		return utility.Quote("")
	case jdb.TypeDataTime:
		return utility.Quote("NOW()")
	default:
		return utility.Quote("")
	}
}

func (s *Postgres) ddlIndex(model *jdb.Model) string {
	var result string
	for _, index := range model.Indices {
		def := ddlIndex(index.Column)

		result = strs.Append(result, def, "\n")
	}

	return result
}

func (s *Postgres) ddlUniqueIndex(model *jdb.Model) string {
	var result string
	for _, index := range model.Uniques {
		def := ddlUniqueIndex(index.Column)

		result = strs.Append(result, def, "\n")
	}

	return result
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

func (s *Postgres) ddlTable(model *jdb.Model) string {
	var columnsDef string
	for _, column := range model.Columns {
		if column.TypeColumn == jdb.TpColumn {
			columnsDef += strs.Format("\n\t%s %s DEFAULT %v,", column.Name, s.typeColumn(column.TypeData), s.defaultValue(column.TypeData))
		}
	}
	columnsDef = strs.Append(columnsDef, ddlPrimaryKey(model), "\n\t")
	result := strs.Format("\nCREATE TABLE IF NOT EXISTS %s (%s\n);", model.Table, columnsDef)
	result = strs.Append(result, s.ddlIndex(model), "\n")
	result = strs.Append(result, s.ddlUniqueIndex(model), "\n\n")
	result = strs.Append(result, s.ddlTriggers(model), "\n\n")

	return result
}

func (s *Postgres) ddlIndexFunction(model *jdb.Model) string {
	result := "\n"
	result = strs.Append(result, s.ddlIndex(model), "\n")
	result = strs.Append(result, s.ddlUniqueIndex(model), "\n")
	result = strs.Append(result, s.ddlTriggers(model), "\n\n")

	return result
}
