package postgres

import (
	"database/sql"
	"fmt"

	"github.com/cgalvisleon/et/logs"
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

/**
* CreateSchema
* @param db *sql.DB, name string
* @return error
**/
func (s *Postgres) CreateSchema(db *sql.DB, name string) error {
	exist, err := existSchema(db, name)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	sql := fmt.Sprintf(`CREATE SCHEMA %s;`, name)
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	logs.Logf(driver, `Schema %s created`, name)

	return nil
}

/**
* DropSchema
* @param db *sql.DB, name string
* @return error
**/
func (s *Postgres) DropSchema(db *sql.DB, name string) error {
	exist, err := existSchema(db, name)
	if err != nil {
		return err
	}

	if !exist {
		return nil
	}

	sql := fmt.Sprintf(`DROP SCHEMA %s;`, name)
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	logs.Logf(driver, `Schema %s droped`, name)

	return nil
}
