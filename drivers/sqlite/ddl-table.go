package sqlite

import (
	"fmt"
	"slices"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

/**
* typeData
* @param tp jdb.TypeData
* @return interface{}
**/
func (s *SqlLite) typeData(tp jdb.TypeData) interface{} {
	switch tp {
	case jdb.TypeDataText:
		return "TEXT"
	case jdb.TypeDataMemo:
		return "TEXT"
	case jdb.TypeDataShortText:
		return "TEXT"
	case jdb.TypeDataKey:
		return "TEXT"
	case jdb.TypeDataNumber:
		return "REAL"
	case jdb.TypeDataInt:
		return "INTEGER"
	case jdb.TypeDataPrecision:
		return "REAL"
	case jdb.TypeDataDateTime:
		return "TEXT"
	case jdb.TypeDataCheckbox:
		return "INTEGER"
	case jdb.TypeDataBytes:
		return "BLOB"
	case jdb.TypeDataObject:
		return "TEXT"
	case jdb.TypeDataSelect:
		return "TEXT"
	case jdb.TypeDataMultiSelect:
		return "TEXT"
	case jdb.TypeDataGeometry:
		return "TEXT"
	case jdb.TypeDataFullText:
		return "TEXT"
	case jdb.TypeDataState:
		return "TEXT"
	case jdb.TypeDataUser:
		return "TEXT"
	case jdb.TypeDataFilesMedia:
		return "TEXT"
	case jdb.TypeDataUrl:
		return "TEXT"
	case jdb.TypeDataEmail:
		return "TEXT"
	case jdb.TypeDataPhone:
		return "TEXT"
	case jdb.TypeDataAddress:
		return "TEXT"
	case jdb.TypeDataRelation:
		return "TEXT"
	case jdb.TypeDataRollup:
		return "TEXT"
	default:
		return "TEXT"
	}
}

/**
* strToTypeData
* @param tp string
* @param lenght int
* @return jdb.TypeData
**/
func (s *SqlLite) strToTypeData(tp string, lenght int) jdb.TypeData {
	tp = strs.Uppcase(tp)
	switch tp {
	case "INTEGER":
		return jdb.TypeDataInt
	case "REAL":
		return jdb.TypeDataNumber
	case "TEXT":
		switch lenght {
		case 80:
			return jdb.TypeDataShortText
		case 20:
			return jdb.TypeDataState
		default:
			return jdb.TypeDataMemo
		}
	case "BLOB":
		return jdb.TypeDataBytes
	case "NUMERIC":
		return jdb.TypeDataNumber
	default:
		return jdb.TypeDataText
	}
}

/**
* defaultValue
* @param tp jdb.TypeData
* @return interface{}
**/
func (s *SqlLite) defaultValue(tp jdb.TypeData) interface{} {
	switch tp {
	case jdb.TypeDataNumber:
		return 0.0
	case jdb.TypeDataInt:
		return 0
	case jdb.TypeDataPrecision:
		return 0.0
	case jdb.TypeDataDateTime:
		return "(datetime('now'))"
	case jdb.TypeDataCheckbox:
		return 0
	case jdb.TypeDataBytes:
		return jdb.Quote("")
	case jdb.TypeDataObject:
		return jdb.Quote(et.Json{})
	case jdb.TypeDataSelect:
		return jdb.Quote("")
	case jdb.TypeDataMultiSelect:
		return jdb.Quote([]et.Json{})
	case jdb.TypeDataGeometry:
		return jdb.Quote(et.Json{
			"type":        "Point",
			"coordinates": []float64{0, 0},
		})
	case jdb.TypeDataFullText:
		return jdb.Quote("")
	case jdb.TypeDataState:
		return jdb.Quote(jdb.ACTIVE)
	case jdb.TypeDataUser:
		return jdb.Quote("")
	case jdb.TypeDataFilesMedia:
		return jdb.Quote("")
	case jdb.TypeDataUrl:
		return jdb.Quote("")
	case jdb.TypeDataEmail:
		return jdb.Quote("")
	case jdb.TypeDataPhone:
		return jdb.Quote("")
	case jdb.TypeDataAddress:
		return jdb.Quote("")
	case jdb.TypeDataRelation:
		return jdb.Quote("")
	case jdb.TypeDataRollup:
		return jdb.Quote("")
	default:
		return jdb.Quote("")
	}
}

/**
* ddlTable
* @param model *jdb.Model
* @return string
**/
func (s *SqlLite) ddlTable(model *jdb.Model) string {
	var columnsDef string
	for _, column := range model.Columns {
		if slices.Contains([]*jdb.Column{model.SystemKeyField}, column) {
			def := fmt.Sprintf("\n\t%s %s DEFAULT %v", column.Name, s.typeData(column.TypeData), s.defaultValue(column.TypeData))
			columnsDef = strs.Append(columnsDef, def, ",")
		} else if slices.Contains([]*jdb.Column{model.FullTextField}, column) && column.FullText != nil {
			def := fmt.Sprintf("\n\t%s TEXT", column.Name)
			columnsDef = strs.Append(columnsDef, def, ",")
		} else if column.TypeColumn == jdb.TpColumn {
			def := fmt.Sprintf("\n\t%s %s DEFAULT %v", column.Name, s.typeData(column.TypeData), s.defaultValue(column.TypeData))
			columnsDef = strs.Append(columnsDef, def, ",")
		}
	}

	ddlPrimaryKey := s.ddlPrimaryKey(model)
	columnsDef = strs.Append(columnsDef, ddlPrimaryKey, ",\n\t")
	ddlForeignKeys := s.ddlForeignKeys(model)
	columnsDef = strs.Append(columnsDef, ddlForeignKeys, ",\n\t")
	result := fmt.Sprintf("\nCREATE TABLE IF NOT EXISTS %s (%s\n);", tableName(model), columnsDef)

	return result
}

/**
* ddlTableRename
* @param oldName string
* @param newName string
* @return string
**/
func (s *SqlLite) ddlTableRename(oldName, newName string) string {
	result := fmt.Sprintf(`ALTER TABLE %s RENAME TO %s;`, oldName, newName)
	return result
}

/**
* ddlTableInsertTo
* @param model *jdb.Model
* @param tableOrigin string
* @return string
**/
func (s *SqlLite) ddlTableInsertTo(model *jdb.Model, tableOrigin string) string {
	fields := ""
	for _, column := range model.Columns {
		if column.TypeColumn == jdb.TpColumn {
			fields = strs.Append(fields, fmt.Sprintf("%s", column.Name), ", ")
		}
	}
	result := fmt.Sprintf("INSERT INTO %s (%s) SELECT %s FROM %s;", tableName(model), fields, fields, tableOrigin)
	return result
}

/**
* ddlTableDrop
* @param table string
* @return string
**/
func (s *SqlLite) ddlTableDrop(table string) string {
	result := fmt.Sprintf("DROP TABLE IF EXISTS %s;", table)
	return result
}

/**
* ddlTableEmpty
* @param table string
* @return string
**/
func (s *SqlLite) ddlTableEmpty(table string) string {
	result := fmt.Sprintf("DELETE FROM %s;", table)
	return result
}
