package postgres

import (
	"fmt"
	"slices"
	"strings"

	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

/**
* ddlIndex
* @param name string
* @param col *jdb.Column
* @return string
**/
func ddlIndex(name string, col *jdb.Column) string {
	result := ""
	if slices.Contains([]jdb.TypeData{jdb.TypeDataObject}, col.TypeData) {
		result = jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $1 ON $2 USING GIN($3 jsonb_path_ops);`, name, col.Model.Table, col.Name)
	} else if slices.Contains([]jdb.TypeData{jdb.TypeDataFullText}, col.TypeData) {
		result = jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $1 ON $2 USING GIN($3);`, name, col.Model.Table, col.Name)
	} else if col.TypeColumn == jdb.TpAtribute && col.Model.SourceField != nil {
		result = jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $1 ON $2 (($3->>'$4'));`, name, col.Model.Table, col.Model.SourceField.Name, col.Name)
	} else {
		result = jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $1 ON $2($3);`, name, col.Model.Table, col.Name)
	}

	return result
}

/**
* ddlUniqueIndex
* @param name string
* @param col *jdb.Column
* @return string
**/
func ddlUniqueIndex(name string, col *jdb.Column) string {
	result := ""
	if col.TypeColumn == jdb.TpColumn {
		result = jdb.SQLDDL(`CREATE UNIQUE INDEX IF NOT EXISTS $1 ON $2($3);`, name, col.Model.Table, col.Name)
	}

	return result
}

/**
* ddlPrimaryKey
* @param model *jdb.Model
* @return string
**/
func (s *Postgres) ddlPrimaryKey(model *jdb.Model) string {
	var result string
	primaryKeys := func() []string {
		var result []string
		for _, v := range model.PrimaryKeys {
			result = append(result, v.Name)
		}

		return result
	}

	if len(primaryKeys()) > 0 {
		result = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s_pk PRIMARY KEY (%s);", model.Table, model.Name, strings.Join(primaryKeys(), ", "))
	}

	return result
}

/**
* ddlForeignKeys
* @param model *jdb.Model
* @return string
**/
func (s *Postgres) ddlForeignKeys(model *jdb.Model) string {
	var result string
	for name, relation := range model.ForeignKeys {
		reference := relation.With
		if reference == nil {
			continue
		}

		referenceKey := ""
		key := ""
		for fkn, pkn := range relation.Fk {
			key = strs.Append(key, fkn, ", ")
			referenceKey = strs.Append(referenceKey, pkn, ", ")
		}
		def := fmt.Sprintf(`ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s)`, model.Table, name, key, reference.Table, referenceKey)
		if relation.OnDeleteCascade {
			def = def + " ON DELETE CASCADE"
		}
		if relation.OnUpdateCascade {
			def = def + " ON UPDATE CASCADE"
		}
		def = def + ";"
		result = strs.Append(result, def, "\n")
	}

	return result
}

/**
* ddlIndex
* @param model *jdb.Model
* @return string
**/
func (s *Postgres) ddlIndex(model *jdb.Model) string {
	var result string
	for name, index := range model.Indices {
		def := ""
		if index.Column.TypeColumn == jdb.TpAtribute && s.version >= 13 {
			def = ddlIndex(name, index.Column)
		} else if index.Column.TypeColumn == jdb.TpColumn {
			def = ddlIndex(name, index.Column)
		}

		result = strs.Append(result, def, "\n")
	}

	return result
}

/**
* ddlUniqueIndex
* @param model *jdb.Model
* @return string
**/
func (s *Postgres) ddlUniqueIndex(model *jdb.Model) string {
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
func (s *Postgres) ddlTableIndex(model *jdb.Model) string {
	result := ""
	result = strs.Append(result, s.ddlIndex(model), "\n")
	result = strs.Append(result, s.ddlPrimaryKey(model), "\n")
	result = strs.Append(result, s.ddlForeignKeys(model), "\n")
	result = strs.Append(result, s.ddlUniqueIndex(model), "\n")

	return fmt.Sprintf("\n%s", result)
}
