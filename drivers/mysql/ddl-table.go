package mysql

import (
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* typeData
* @param tp jdb.TypeData
* @return interface{}
**/
func (s *Mysql) typeData(tp jdb.TypeData) interface{} {
	switch tp {
	case jdb.TypeDataText:
		return "VARCHAR(250)"
	case jdb.TypeDataMemo:
		return "TEXT"
	case jdb.TypeDataShortText:
		return "VARCHAR(80)"
	case jdb.TypeDataKey:
		return "VARCHAR(80)"
	case jdb.TypeDataNumber:
		return "DECIMAL(18,2)"
	case jdb.TypeDataInt:
		return "BIGINT"
	case jdb.TypeDataPrecision:
		return "DOUBLE PRECISION"
	case jdb.TypeDataDateTime:
		return "TIMESTAMP"
	case jdb.TypeDataCheckbox:
		return "BOOLEAN"
	case jdb.TypeDataBytes:
		return "BYTEA"
	case jdb.TypeDataObject:
		return "JSONB"
	case jdb.TypeDataSelect:
		return "VARCHAR(250)"
	case jdb.TypeDataMultiSelect:
		return "JSONB"
	case jdb.TypeDataGeometry:
		return "JSONB"
	case jdb.TypeDataFullText:
		return "TSVECTOR"
	case jdb.TypeDataState:
		return "VARCHAR(80)"
	case jdb.TypeDataUser:
		return "VARCHAR(250)"
	case jdb.TypeDataFilesMedia:
		return "TEXT"
	case jdb.TypeDataUrl:
		return "TEXT"
	case jdb.TypeDataEmail:
		return "VARCHAR(250)"
	case jdb.TypeDataPhone:
		return "VARCHAR(250)"
	case jdb.TypeDataAddress:
		return "TEXT"
	case jdb.TypeDataRelation:
		return "VARCHAR(80)"
	case jdb.TypeDataRollup:
		return "VARCHAR(80)"
	default:
		return "VARCHAR(250)"
	}
}

/**
* strToTypeData
* @param tp string
* @param lenght int
* @return jdb.TypeData
**/
func (s *Mysql) strToTypeData(tp string, lenght int) jdb.TypeData {
	tp = strs.Uppcase(tp)
	switch tp {
	case "BOOLEAN":
		return jdb.TypeDataCheckbox
	case "INTEGER":
		return jdb.TypeDataInt
	case "INT4":
		return jdb.TypeDataInt
	case "VARCHAR":
		switch lenght {
		case 80:
			return jdb.TypeDataShortText
		case 20:
			return jdb.TypeDataShortText
		default:
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
		return jdb.TypeDataInt
	case "VARCHAR(250)":
		return jdb.TypeDataText
	case "TIMESTAMP":
		return jdb.TypeDataDateTime
	case "BYTEA":
		return jdb.TypeDataBytes
	case "TSVECTOR":
		return jdb.TypeDataFullText
	default:
		return jdb.TypeDataText
	}
}

/**
* defaultValue
* @param tp jdb.TypeData
* @return interface{}
**/
func (s *Mysql) defaultValue(tp jdb.TypeData) interface{} {
	switch tp {
	case jdb.TypeDataNumber:
		return 0.0
	case jdb.TypeDataInt:
		return 0
	case jdb.TypeDataPrecision:
		return 0.0
	case jdb.TypeDataDateTime:
		return utility.Quote("NOW()")
	case jdb.TypeDataCheckbox:
		return utility.Quote(false)
	case jdb.TypeDataBytes:
		return utility.Quote("")
	case jdb.TypeDataObject:
		return utility.Quote(et.Json{})
	case jdb.TypeDataSelect:
		return utility.Quote("")
	case jdb.TypeDataMultiSelect:
		return utility.Quote([]et.Json{})
	case jdb.TypeDataGeometry:
		return utility.Quote(et.Json{
			"type":        "Point",
			"coordinates": []float64{0, 0},
		})
	case jdb.TypeDataFullText:
		return utility.Quote("")
	case jdb.TypeDataState:
		return utility.Quote(utility.ACTIVE)
	case jdb.TypeDataUser:
		return utility.Quote("")
	case jdb.TypeDataFilesMedia:
		return utility.Quote("")
	case jdb.TypeDataUrl:
		return utility.Quote("")
	case jdb.TypeDataEmail:
		return utility.Quote("")
	case jdb.TypeDataPhone:
		return utility.Quote("")
	case jdb.TypeDataAddress:
		return utility.Quote("")
	case jdb.TypeDataRelation:
		return utility.Quote("")
	case jdb.TypeDataRollup:
		return utility.Quote("")
	default:
		return utility.Quote("")
	}
}

/**
* ddlTable
* @param model *jdb.Model
* @return string
**/
func (s *Mysql) ddlTable(model *jdb.Model) string {
	var columnsDef string
	for _, column := range model.Columns {
		if slices.Contains([]*jdb.Column{model.SystemKeyField}, column) {
			if s.version >= 15 {
				def := strs.Format("\n\t%s %s DEFAULT %v INVISIBLE", column.Name, s.typeData(column.TypeData), s.defaultValue(column.TypeData))
				columnsDef = strs.Append(columnsDef, def, ",")
			} else {
				def := strs.Format("\n\t%s %s DEFAULT %v", column.Name, s.typeData(column.TypeData), s.defaultValue(column.TypeData))
				columnsDef = strs.Append(columnsDef, def, ",")
			}
		} else if slices.Contains([]*jdb.Column{model.FullTextField}, column) && column.FullText != nil {
			columns := ""
			for _, col := range column.FullText.Columns {
				columns = strs.Append(columns, strs.Format("COALESCE(%s, '')", col), " || ' ' || ")
			}
			def := strs.Format("\n\t%s TSVECTOR GENERATED ALWAYS AS (to_tsvector('%s', %s)) STORED", column.Name, column.FullText.Language, columns)
			columnsDef = strs.Append(columnsDef, def, ",")
		} else if column.TypeColumn == jdb.TpColumn {
			def := strs.Format("\n\t%s %s DEFAULT %v", column.Name, s.typeData(column.TypeData), s.defaultValue(column.TypeData))
			columnsDef = strs.Append(columnsDef, def, ",")
		}
	}
	result := strs.Format("\nCREATE TABLE IF NOT EXISTS %s (%s\n);", tableName(model), columnsDef)

	return result
}

/**
* ddlTableRename
* @param oldName string
* @param newName string
* @return string
**/
func (s *Mysql) ddlTableRename(oldName, newName string) string {
	result := strs.Format(`ALTER TABLE %s RENAME TO %s;`, oldName, newName)

	return result
}

/**
* ddlTableInsertTo
* @param model *jdb.Model
* @param tableOrigin string
* @return string
**/
func (s *Mysql) ddlTableInsertTo(model *jdb.Model, tableOrigin string) string {
	fields := ""
	for _, column := range model.Columns {
		if column.TypeColumn == jdb.TpColumn {
			fields = strs.Append(fields, strs.Format("%s", column.Name), ", ")
		}
	}
	result := strs.Format("INSERT INTO %s (%s)\nSELECT %s FROM %s;", tableName(model), fields, fields, tableOrigin)

	return result
}

/**
* ddlTableDrop
* @param table string
* @return string
**/
func (s *Mysql) ddlTableDrop(table string) string {
	result := strs.Format("DROP TABLE IF EXISTS %s CASCADE;", table)

	return result
}
