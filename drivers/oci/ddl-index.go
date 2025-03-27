package oci

import (
	"slices"
	"strings"

	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

func ddlIndex(name string, col *jdb.Column) string {
	result := ""
	if slices.Contains([]jdb.TypeData{jdb.TypeDataObject, jdb.TypeDataArray}, col.TypeData) {
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

func ddlUniqueIndex(name string, col *jdb.Column) string {
	result := ""
	if col.TypeColumn == jdb.TpColumn {
		result = jdb.SQLDDL(`CREATE UNIQUE INDEX IF NOT EXISTS $1 ON $2($3);`, name, col.Model.Table, col.Name)
	}

	return result
}

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
		result = strs.Format("PRIMARY KEY (%s)", strings.Join(primaryKeys(), ", "))
	}

	return result
}

func (s *Postgres) ddlForeignKeys(model *jdb.Model) string {
	var result string
	for key, fk := range model.ForeignKeys {
		ref := fk.Detail.Fk
		def := strs.Format(`ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s)`, model.Table, key, fk.Name, ref.Model.Table, ref.Name)
		if fk.Detail.OnDeleteCascade {
			def = def + " ON DELETE CASCADE"
		}
		if fk.Detail.OnUpdateCascade {
			def = def + " ON UPDATE CASCADE"
		}
		def = strs.Format("SELECT core.add_constraint_if_not_exists('%s', '%s', '%s', '%s');", model.Schema.Low(), model.Low(), key, def)
		result = strs.Append(result, def, "\n")
	}

	return result
}

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

func (s *Postgres) ddlIndexFunction(model *jdb.Model) string {
	result := ""
	result = strs.Append(result, s.ddlIndex(model), "\n")
	result = strs.Append(result, s.ddlUniqueIndex(model), "\n")
	result = strs.Append(result, s.ddlForeignKeys(model), "\n\n")
	result = strs.Append(result, s.ddlTriggers(model), "\n\n")

	return result
}

func (s *Postgres) ddlSystemFunction(model *jdb.Model) string {
	result := ""
	result = strs.Append(result, s.ddlTriggers(model), "\n\n")

	return result
}
