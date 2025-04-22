package sqlite

import (
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
)

func parceSQL(sql string) string {
	return strs.Change(sql,
		[]string{"date_make", "date_update", "_id", "_idt", "_state"},
		[]string{jdb.CREATED_AT, jdb.UPDATED_AT, jdb.PRIMARYKEY, jdb.SYSID, jdb.STATUS})
}
