package mysql

import (
	"database/sql"
	"errors"

	"github.com/cgalvisleon/et/console"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* connectTo
* @param connStr string
* @return *sql.DB, error
**/
func (s *Mysql) connectTo(chain string) (*sql.DB, error) {
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
* exeistDatabase
* @param db *sql.DB, name string
* @return bool, error
**/
func (s *Mysql) ExistDatabase(db *sql.DB, name string) (bool, error) {
	sql := jdb.SQLDDL(`
	SELECT SCHEMA_NAME 
	FROM INFORMATION_SCHEMA.SCHEMATA 
	WHERE SCHEMA_NAME = $1`, name)
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
func (s *Mysql) CreateDatabase(db *sql.DB, name string) error {
	if s.jdb == nil {
		return errors.New(MSG_JDB_NOT_DEFINED)
	}

	exist, err := s.ExistDatabase(db, name)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	sql := `CREATE DATABASE $1;`
	_, err = db.Exec(sql, name)
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
func (s *Mysql) DropDatabase(db *sql.DB, name string) error {
	if s.jdb == nil {
		return errors.New(MSG_JDB_NOT_DEFINED)
	}

	exist, err := s.ExistDatabase(db, name)
	if err != nil {
		return err
	}

	if !exist {
		return nil
	}

	sql := `DROP DATABASE $1;`
	_, err = db.Exec(sql, name)
	if err != nil {
		return err
	}

	console.Logf(s.name, `Database %s droped`, name)

	return nil
}

/**
* Connect
* @param connection jdb.ConnectParams
* @return error
**/
func (s *Mysql) Connect(connection jdb.ConnectParams) (*sql.DB, error) {
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
