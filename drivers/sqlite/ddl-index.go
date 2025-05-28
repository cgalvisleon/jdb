package sqlite

import (
	"strings"

	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

/**
* ddlUniqueIndex
* @param name string
* @param col *jdb.Column
* @return string
**/
func ddlUniqueIndex(name string, col *jdb.Column) string {
	result := ""
	if col.TypeColumn == jdb.TpColumn {
		result = jdb.SQLDDL(`CREATE UNIQUE INDEX IF NOT EXISTS $1 ON $2($3);`, name, tableName(col.Model), col.Name)
	}

	return result
}

/**
* ddlPrimaryKey
* @param model *jdb.Model
* @return string
**/
func (s *SqlLite) ddlPrimaryKey(model *jdb.Model) string {
	var result string
	primaryKeys := func() []string {
		var result []string
		for _, v := range model.PrimaryKeys {
			result = append(result, v.Name)
		}
		return result
	}

	if len(primaryKeys()) > 0 {
		result = strs.Format("PRIMARY KEY (%s)", strings.Join(primaryKeys(), ", "))
	}
	return result
}

/**
* ddlForeignKeys
* @param model *jdb.Model
* @return string
**/
func (s *SqlLite) ddlForeignKeys(model *jdb.Model) string {
	var result string
	for _, column := range model.Columns {
		if column.TypeColumn == jdb.TpColumn {
			if column.Detail != nil && column.Detail.With != nil {
				with := column.Detail.With
				fk := column.Detail.Fk
				if len(fk) > 0 {
					fks := ""
					references := ""
					for fkn, pkn := range fk {
						fks = strs.Append(fks, strs.Format("%s", fkn), ", ")
						references = strs.Append(references, strs.Format("%s", pkn), ", ")
					}
					fks = strings.TrimSuffix(fks, ", ")
					references = strings.TrimSuffix(references, ", ")

					onDelete := ""
					if column.Detail.OnDeleteCascade {
						onDelete = " ON DELETE CASCADE"
					}

					onUpdate := ""
					if column.Detail.OnUpdateCascade {
						onUpdate = " ON UPDATE CASCADE"
					}

					idx := strs.Format("FOREIGN KEY (%s) REFERENCES %s (%s)%s%s",
						fks,
						tableName(with),
						references,
						onDelete,
						onUpdate)
					result = strs.Append(result, idx, "\n")
				}
			}
		}
	}

	return result
}

/**
* ddlIndex
* @param model *jdb.Model
* @return string
**/
func (s *SqlLite) ddlIndex(model *jdb.Model) string {
	var result string
	for _, column := range model.Columns {
		if column.TypeColumn == jdb.TpColumn {
			if column.IsKeyfield {
				idx := strs.Format("CREATE INDEX IF NOT EXISTS idx_%s_%s ON %s (%s);", model.Name, column.Name, tableName(model), column.Name)
				result = strs.Append(result, idx, "\n")
			}
		}
	}

	return result
}

/**
* ddlUniqueIndex
* @param model *jdb.Model
* @return string
**/
func (s *SqlLite) ddlUniqueIndex(model *jdb.Model) string {
	var result string
	for name, index := range model.Uniques {
		def := ""
		if index.Column.TypeColumn == jdb.TpColumn {
			def = ddlUniqueIndex(name, index.Column)
		}

		result = strs.Append(result, def, "\n")
	}

	return result
}

/**
* ddlTableIndex
* @param model *jdb.Model
* @return string
**/
func (s *SqlLite) ddlTableIndex(model *jdb.Model) string {
	result := ""
	result = strs.Append(result, s.ddlIndex(model), "\n")
	result = strs.Append(result, s.ddlUniqueIndex(model), "\n")

	return strs.Format("\n%s", result)
}
