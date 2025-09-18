package postgres

import (
	"database/sql"
	"fmt"

	"github.com/cgalvisleon/et/console"
	jdb "github.com/cgalvisleon/jdb/congo"
	_ "github.com/lib/pq"
)

/**
* Connect
* @param connection jdb.ConnectParams
* @return *sql.DB, error
**/
func (s *Postgres) Connect(database *jdb.Database) (*sql.DB, error) {
	s.database = database
	s.name = database.Name

	defaultChain, err := s.connection.defaultChain()
	if err != nil {
		return nil, err
	}

	db, err := s.connectTo(defaultChain)
	if err != nil {
		return nil, err
	}

	err = s.createDatabase(db, database.Name)
	if err != nil {
		return nil, err
	}

	if db != nil {
		err := db.Close()
		if err != nil {
			return nil, err
		}
	}

	s.connection.Database = database.Name
	chain, err := s.connection.chain()
	if err != nil {
		return nil, err
	}

	db, err = s.connectTo(chain)
	if err != nil {
		return nil, err
	}

	console.Logf(driver, `Connected to %s:%s`, s.connection.Host, s.connection.Database)

	return db, nil
}

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

	console.Logf(driver, `Database %s created`, name)

	return nil
}

/**
* DropDatabase
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

	console.Logf(driver, `Database %s droped`, name)

	return nil
}
