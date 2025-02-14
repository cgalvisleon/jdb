package postgres

import (
	"slices"

	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) typeData(tp jdb.TypeData) interface{} {
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
	case jdb.TypeDataBytes:
		return "BYTEA"
	case jdb.TypeDataGeometry:
		return "JSONB"
	case jdb.TypeDataFullText:
		return "TSVECTOR"
	default:
		return "VARCHAR(250)"
	}
}

func (s *Postgres) strToTypeData(tp string, lenght int) jdb.TypeData {
	tp = strs.Uppcase(tp)
	switch tp {
	case "ARRAY":
		return jdb.TypeDataArray
	case "BOOLEAN":
		return jdb.TypeDataBool
	case "INTEGER":
		return jdb.TypeDataInt
	case "INT4":
		return jdb.TypeDataInt
	case "VARCHAR":
		if lenght == 80 {
			return jdb.TypeDataKey
		} else if lenght == 20 {
			return jdb.TypeDataShortText
		} else {
			return jdb.TypeDataText
		}
	case "VARCHAR(80)":
		return jdb.TypeDataKey
	case "VARCHAR(20)":
		return jdb.TypeDataState
	case "TEXT":
		return jdb.TypeDataMemo
	case "DECIMAL(18,2)":
		return jdb.TypeDataNumber
	case "DOUBLE PRECISION":
		return jdb.TypeDataPrecision
	case "NUMERIC":
		return jdb.TypeDataNumber
	case "JSONB":
		return jdb.TypeDataObject
	case "BIGINT":
		return jdb.TypeDataSerie
	case "VARCHAR(250)":
		return jdb.TypeDataText
	case "TIMESTAMP":
		return jdb.TypeDataTime
	case "BYTEA":
		return jdb.TypeDataBytes
	case "TSVECTOR":
		return jdb.TypeDataFullText
	default:
		return jdb.TypeDataText
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
	case jdb.TypeDataBytes:
		return utility.Quote("")
	case jdb.TypeDataGeometry:
		return utility.Quote("{type: 'Point', coordinates: [0, 0]}")
	default:
		return utility.Quote("")
	}
}

func (s *Postgres) ddlTable(model *jdb.Model) string {
	var columnsDef string
	for _, column := range model.Columns {
		if slices.Contains([]*jdb.Column{model.SystemKeyField, model.FullTextField}, column) {
			if s.version >= 13 {
				def := strs.Format("\n\t%s %s INVISIBLE DEFAULT %v", column.Name, s.typeData(column.TypeData), s.defaultValue(column.TypeData))
				columnsDef = strs.Append(columnsDef, def, ",")
			} else {
				def := strs.Format("\n\t%s %s DEFAULT %v", column.Name, s.typeData(column.TypeData), s.defaultValue(column.TypeData))
				columnsDef = strs.Append(columnsDef, def, ",")
			}
		} else if column.TypeColumn == jdb.TpColumn {
			def := strs.Format("\n\t%s %s DEFAULT %v", column.Name, s.typeData(column.TypeData), s.defaultValue(column.TypeData))
			columnsDef = strs.Append(columnsDef, def, ",")
		}
	}
	def := s.ddlPrimaryKey(model)
	columnsDef = strs.Append(columnsDef, def, ",\n\t")
	result := strs.Format("\nCREATE TABLE IF NOT EXISTS %s (%s\n);", model.Table, columnsDef)
	result = strs.Append(result, s.ddlIndexFunction(model), "\n")

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
			fields = strs.Append(fields, strs.Format("%s", column.Name), ", ")
		}
	}
	result := strs.Format("INSERT INTO %s (%s)\nSELECT %s FROM %s;", old.Table, fields, fields, backupTable)

	return result
}

func (s *Postgres) ddlTableDrop(table string) string {
	result := strs.Format("DROP TABLE IF EXISTS %s CASCADE;", table)

	return result
}

func (s *Postgres) tableExists(schema, tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = $1
			AND table_name = $2
		);`

	var exists bool
	err := s.db.QueryRow(query, schema, tableName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
