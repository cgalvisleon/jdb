package postgres

import (
	"database/sql"
	"fmt"

	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/jdb/jdb"
	_ "github.com/lib/pq"
)

/**
* connectTo
* @param chain string
* @return *sql.DB, error
**/
func (s *Postgres) connectTo(chain string) (*sql.DB, error) {
	db, err := sql.Open(driver, chain)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

/**
* existDatabase
* @param db *DB, name string
* @return bool, error
**/
func (s *Postgres) existDatabase(db *sql.DB, name string) (bool, error) {
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
func (s *Postgres) createDatabase(db *sql.DB, name string) error {
	exist, err := s.existDatabase(db, name)
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
	exist, err := s.existDatabase(db, name)
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
