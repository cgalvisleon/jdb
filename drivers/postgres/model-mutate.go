package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) ddlMutate(old, model *jdb.Model) string {
	backupTable := strs.Format(`%s_BACKUP`, old.Up())
	result := "\n"
	result = strs.Append(result, s.ddlTableRename(old.Table, backupTable), "\n")
	result = strs.Append(result, s.ddlTable(model), "\n")
	result = strs.Append(result, s.ddlTableInsert(old), "\n\n")
	result = strs.Append(result, s.ddlTableDrop(model.Schema.Name+"."+backupTable), "\n\n")

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
			fields = strs.Append(fields, strs.Format("%s", column.Up()), ", ")
		}
	}
	result := strs.Format("INSERT INTO %s (%s)\nSELECT %s FROM %s;", old.Table, fields, fields, backupTable)

	return result
}

func (s *Postgres) ddlTableDrop(table string) string {
	result := strs.Format("DROP TABLE IF EXISTS %s CASCADE;", table)

	return result
}
