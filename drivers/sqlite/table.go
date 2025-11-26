package sqlite

import "database/sql"

func ExistTable(db *sql.DB, schema, table string) (bool, error) {
	return false, nil
}
