package oci

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Oracle) ddlMutate(old, model *jdb.Model, drop bool) string {
	backupTable := strs.Format(`%s_backup`, old.Up())
	result := "\n"
	result = strs.Append(result, s.ddlTableRename(old.Table, backupTable), "\n")
	result = strs.Append(result, s.ddlTable(model), "\n")
	result = strs.Append(result, s.ddlTableInsert(old), "\n\n")
	if drop {
		result = strs.Append(result, s.ddlTableDrop(model.Schema.Name+"."+backupTable), "\n\n")
	}

	return result
}
