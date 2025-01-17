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
	case jdb.TypeDataState:
		return "VARCHAR(20)"
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
	case jdb.TypeDataGeometry:
		return "JSONB"
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
	case jdb.TypeDataState:
		return utility.Quote(utility.ACTIVE)
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
	case jdb.TypeDataGeometry:
		return utility.Quote("{type: 'Point', coordinates: [0, 0]}")
	default:
		return utility.Quote("")
	}
}

func (s *Postgres) ddlTable(model *jdb.Model) string {
	var columnsDef string
	for _, column := range model.Columns {
		if column == model.SystemKeyField {
			columnsDef += strs.Format("\n\t%s %s INVISIBLE DEFAULT %v,", column.Name, s.typeColumn(column.TypeData), s.defaultValue(column.TypeData))
		} else if column.TypeColumn == jdb.TpColumn {
			columnsDef += strs.Format("\n\t%s %s DEFAULT %v,", column.Name, s.typeColumn(column.TypeData), s.defaultValue(column.TypeData))
		}
	}
	columnsDef = strs.Append(columnsDef, s.ddlPrimaryKey(model), "\n\t")
	result := strs.Format("\nCREATE TABLE IF NOT EXISTS %s (%s\n);", model.Table, columnsDef)
	result = strs.Append(result, s.ddlIndex(model), "\n")
	result = strs.Append(result, s.ddlUniqueIndex(model), "\n\n")
	result = strs.Append(result, s.ddlForeignKeys(model), "\n\n")
	result = strs.Append(result, s.ddlTriggers(model), "\n\n")

	return result
}

func (s *Postgres) ddlTableRename(old, new string) string {
	result := strs.Format(`ALTER TABLE %s RENAME TO %s;`, old, new)

	return result
}

func (s *Postgres) ddlTableInsert(old *jdb.Model) string {
	backupTable := strs.Format(`%s_BACKUP`, old.Table)
	fields := ""
	for _, column := range old.Columns {
		if column.TypeColumn == jdb.TpColumn {
			fields = strs.Append(fields, strs.Format("%s", column.Up()), ", ")
		}
	}
	result := strs.Format("INSERT INTO %s (%s)\nSELECT %s FROM %s;", old.Table, fields, fields, backupTable)

	return result
}

func (s *Postgres) ddlTableDrop(table string) string {
	result := strs.Format("DROP TABLE IF EXISTS %s CASCADE;", table)

	return result
}
