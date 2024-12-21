package postgres

import (
	"slices"
	"strings"

	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

func ddlIndex(col *jdb.Column) string {
	result := jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $2_$3_$4_IDX ON $1($4);`, col.Model.Table, col.Model.Schema.Name, col.Model.Name, col.Field)
	if slices.Contains([]jdb.TypeData{jdb.TypeDataObject, jdb.TypeDataArray}, col.TypeData) {
		result = jdb.SQLDDL(`CREATE INDEX IF NOT EXISTS $2_$3_$4_IDX ON $1 USING GIN($4);`, col.Model.Table, col.Model.Schema.Name, col.Model.Name, col.Field)
	}

	return strs.Uppcase(result)
}

func ddlUniqueIndex(col *jdb.Column) string {
	result := jdb.SQLDDL(`CREATE UNIQUE INDEX IF NOT EXISTS $2_$3_$4_IDX ON $1($4);`, col.Model.Table, col.Model.Schema.Name, col.Model.Name, col.Field)
	if slices.Contains([]jdb.TypeData{jdb.TypeDataObject, jdb.TypeDataArray}, col.TypeData) {
		result = ""
	}

	return strs.Uppcase(result)
}

func ddlPrimaryKey(model *jdb.Model) string {
	primaryKeys := func() []string {
		var result []string
		for _, v := range model.Keys {
			result = append(result, v.Field)
		}

		return result
	}

	result := strs.Format("PRIMARY KEY (%s)", strings.Join(primaryKeys(), ", "))

	return strs.Uppcase(result)
}

func ddlForeignKeys(model *jdb.Model) string {
	var result string
	for _, ref := range model.References {
		field := ref.Key.Field
		key := field + "_FKEY"
		key = strs.Replace(key, "-", "_")
		def := strs.Format(`ALTER TABLE IF EXISTS %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s);`, model.Table, key, field, ref.To.Table, ref.To.Field)
		if ref.OnDeleteCascade {
			def = def + " ON DELETE CASCADE"
		}
		if ref.OnUpdateCascade {
			def = def + " ON UPDATE CASCADE"
		}
		result = strs.Append(result, def, "\n")
	}

	return strs.Uppcase(result)
}
