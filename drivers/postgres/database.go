package postgres

import (
	"database/sql"
	"fmt"

	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/jdb/jdb"
	_ "github.com/lib/pq"
)

/**
* existDatabase
* @param db *DB, name string
* @return bool, error
**/
func existDatabase(db *sql.DB, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
	SELECT 1
	FROM pg_database
	WHERE UPPER(datname) = UPPER($1));`
	rows, err := db.Query(sql, name)
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

/**
* createDatabase
* @param db *sql.DB, name string
* @return error
**/
func (s *Postgres) CreateDatabase(db *sql.DB, name string) error {
	exist, err := existDatabase(db, name)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	sql := fmt.Sprintf(`CREATE DATABASE %s;`, name)
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	logs.Logf(driver, `Database %s created`, name)

	return nil
}

/**
* dropDatabase
* @param db *sql.DB, name string
* @return error
**/
func (s *Postgres) DropDatabase(db *sql.DB, name string) error {
	exist, err := existDatabase(db, name)
	if err != nil {
		return err
	}

	if !exist {
		return nil
	}

	sql := fmt.Sprintf(`DROP DATABASE %s;`, name)
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	logs.Logf(driver, `Database %s droped`, name)

	return nil
}
