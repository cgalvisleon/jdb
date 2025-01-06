package postgres

import (
	"slices"

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
