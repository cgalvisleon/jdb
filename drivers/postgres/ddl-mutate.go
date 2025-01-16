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
