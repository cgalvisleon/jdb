package postgres

import (
	"database/sql"
	"fmt"

	"github.com/cgalvisleon/et/console"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* connectTo
* @param chain string
* @return *sql.DB, error
**/
func (s *Postgres) connectTo(chain string) (*sql.DB, error) {
	db, err := sql.Open(s.name, chain)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

/**
* ExistDatabase
* @param db *DB, name string
* @return bool, error
**/
func (s *Postgres) ExistDatabase(db *sql.DB, name string) (bool, error) {
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
* CreateDatabase
* @param db *sql.DB, name string
* @return error
**/
func (s *Postgres) CreateDatabase(db *sql.DB, name string) error {
	exist, err := s.ExistDatabase(db, name)
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

	console.Logf(s.name, `Database %s created`, name)

	return nil
}

/**
* DropDatabase
* @param db *sql.DB, name string
* @return error
**/
func (s *Postgres) DropDatabase(db *sql.DB, name string) error {
	exist, err := s.ExistDatabase(db, name)
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

	console.Logf(s.name, `Database %s droped`, name)

	return nil
}

/**
* Connect
* @param connection jdb.ConnectParams
* @return *sql.DB, error
**/
func (s *Postgres) Connect(connection jdb.ConnectParams) (*sql.DB, error) {
	if s.jdb == nil {
		return nil, fmt.Errorf(MSG_JDB_NOT_DEFINED)
	}

	defaultChain, err := s.connection.defaultChain()
	if err != nil {
		return nil, err
	}

	db, err := s.connectTo(defaultChain)
	if err != nil {
		return nil, err
	}

	params := connection.Params.(*Connection)
	params.Database = connection.Name
	err = s.CreateDatabase(db, params.Database)
	if err != nil {
		return nil, err
	}

	if db != nil {
		err := db.Close()
		if err != nil {
			return nil, err
		}
	}

	chain, err := params.Chain()
	if err != nil {
		return nil, err
	}

	db, err = s.connectTo(chain)
	if err != nil {
		return nil, err
	}

	s.connected = db != nil
	console.Logf(s.name, `Connected to %s:%s`, params.Host, params.Database)

	return db, nil
}
