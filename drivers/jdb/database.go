package jdb

import "database/sql"

func ExistDatabase(db *sql.DB, name string) (bool, error) {
	return false, nil
}

func CreateDatabase(db *sql.DB, name string) error {
	return nil
}

func DropDatabase(db *sql.DB, name string) error {
	return nil
}
