package postgres

import (
	"slices"
	"strings"

	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

func ddlIndex(col *jdb.Column) string {
	result := ""
	if slices.Contains([]jdb.TypeData{jdb.TypeDataObject, jdb.TypeDataArray}, col.TypeData) {
		result = jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $2_$3_$4_idx ON $1 USING GIN($4 jsonb_path_ops);`, col.Model.Table, col.Model.Schema.Name, col.Model.Name, col.Field)
	} else if slices.Contains([]jdb.TypeData{jdb.TypeDataFullText}, col.TypeData) {
		result = jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $2_$3_$4_idx ON $1 USING GIN($4);`, col.Model.Table, col.Model.Schema.Name, col.Model.Name, col.Field)
	} else if col.TypeColumn == jdb.TpColumn {
		result = jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $2_$3_$4_idx ON $1($4);`, col.Model.Table, col.Model.Schema.Name, col.Model.Name, col.Field)
	}

	return result
}

func ddlUniqueIndex(col *jdb.Column) string {
	result := ""
	if col.TypeColumn == jdb.TpColumn {
		result = jdb.SQLDDL(`CREATE UNIQUE INDEX IF NOT EXISTS $2_$3_$4_idx ON $1($4);`, col.Model.Table, col.Model.Schema.Name, col.Model.Name, col.Field)
	}

	return result
}

func (s *Postgres) ddlPrimaryKey(model *jdb.Model) string {
	primaryKeys := func() []string {
		var result []string
		for _, v := range model.Keys {
			result = append(result, v.Field)
		}

		return result
	}

	result := strs.Format("PRIMARY KEY (%s)", strings.Join(primaryKeys(), ", "))

	return result
}

func (s *Postgres) ddlForeignKeys(model *jdb.Model) string {
	var result string
	for _, ref := range model.References {
		field := ref.Key.Field
		key := field + "_FKEY"
		key = strs.Replace(key, "-", "_")
		key = strs.Uppcase(key)
		def := strs.Format(`ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s);`, model.Table, key, field, ref.To.Table, ref.To.Field)
		if ref.OnDeleteCascade {
			def = def + " ON DELETE CASCADE"
		}
		if ref.OnUpdateCascade {
			def = def + " ON UPDATE CASCADE"
		}
		result = strs.Format("SELECT core.add_constraint_if_not_exists('%s', '%s', '%s', '%s');\n", model.Schema.Low(), model.Low(), strs.Lowcase(key), def)
	}

	return result
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

func (s *Postgres) ddlIndexFunction(model *jdb.Model) string {
	result := "\n"
	result = strs.Append(result, s.ddlIndex(model), "\n")
	result = strs.Append(result, s.ddlUniqueIndex(model), "\n")
	result = strs.Append(result, s.ddlForeignKeys(model), "\n\n")
	result = strs.Append(result, s.ddlTriggers(model), "\n\n")

	return result
}
