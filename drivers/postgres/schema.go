package postgres

import (
	"database/sql"

	"github.com/cgalvisleon/jdb/jdb"
)

/**
* existSchema
* @param db *sql.DB, name string
* @return bool, error
**/
func existSchema(db *sql.DB, name string) (bool, error) {
	rows, err := db.Query(`
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.schemata
		WHERE UPPER(schema_name) = UPPER($1));`, name)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	items := jdb.RowsToItems(rows)

	if items.Count == 0 {
		return false, nil
	}

	return items.Bool(0, "exists"), nil
}
