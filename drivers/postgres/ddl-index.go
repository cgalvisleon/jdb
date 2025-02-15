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
		result = jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $2_$3_$4_idx ON $1 USING GIN($4 jsonb_path_ops);`, col.Model.Table, col.Model.Schema.Name, col.Model.Name, col.Name)
	} else if slices.Contains([]jdb.TypeData{jdb.TypeDataFullText}, col.TypeData) {
		result = jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $2_$3_$4_idx ON $1 USING GIN($4);`, col.Model.Table, col.Model.Schema.Name, col.Model.Name, col.Name)
	} else if col.TypeColumn == jdb.TpAtribute && col.Model.SourceField != nil {
		result = jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $2_$3_$4_idx ON $1 ((%s->>'%s'));`, col.Model.Table, col.Model.Schema.Name, col.Model.Name, col.Model.SourceField.Name, col.Name)
	} else {
		result = jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $2_$3_$4_idx ON $1($4);`, col.Model.Table, col.Model.Schema.Name, col.Model.Name, col.Name)
	}

	return result
}

func ddlUniqueIndex(col *jdb.Column) string {
	result := ""
	if col.TypeColumn == jdb.TpColumn {
		result = jdb.SQLDDL(`CREATE UNIQUE INDEX IF NOT EXISTS $2_$3_$4_idx ON $1($4);`, col.Model.Table, col.Model.Schema.Name, col.Model.Name, col.Name)
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
	for _, fk := range model.ForeignKeys {
		ref := fk.Detail.Fk
		key := ref.Name + "_FKEY"
		key = strs.Replace(key, "-", "_")
		key = strs.Lowcase(key)
		def := strs.Format(`ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s)`, model.Table, key, ref.Name, ref.Model.Table, ref.Name)
		if fk.Detail.OnDeleteCascade {
			def = def + " ON DELETE CASCADE"
		}
		if fk.Detail.OnUpdateCascade {
			def = def + " ON UPDATE CASCADE"
		}
		def = strs.Format("SELECT core.add_constraint_if_not_exists('%s', '%s', '%s', '%s');", model.Schema.Low(), model.Low(), strs.Lowcase(key), def)
		result = strs.Append(result, def, "\n")
	}

	return result
}

func (s *Postgres) ddlIndex(model *jdb.Model) string {
	var result string
	for _, index := range model.Indices {
		def := ""
		if index.Column.TypeColumn == jdb.TpAtribute && s.version >= 13 {
			def = ddlIndex(index.Column)
		} else if index.Column.TypeColumn == jdb.TpColumn {
			def = ddlIndex(index.Column)
		}

		result = strs.Append(result, def, "\n")
	}

	return result
}

func (s *Postgres) ddlUniqueIndex(model *jdb.Model) string {
	var result string
	for _, index := range model.Uniques {
		def := ""
		if index.Column.TypeColumn == jdb.TpColumn {
			def = ddlUniqueIndex(index.Column)
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
