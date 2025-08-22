package mysql

import (
	"database/sql"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/msg"
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
* @param name string
* @return bool, error
**/
func (s *Mysql) ExistDatabase(name string) (bool, error) {
	sql := jdb.SQLDDL(`
	SELECT SCHEMA_NAME 
	FROM INFORMATION_SCHEMA.SCHEMATA 
	WHERE SCHEMA_NAME = $1`, name)
	items, err := jdb.Query(s.db, sql, name)
	if err != nil {
		return false, err
	}

	if items.Count == 0 {
		return false, nil
	}

	return items.Bool(0, "exists"), nil
}

/**
* CreateDatabase
* @param name string
* @return error
**/
func (s *Mysql) CreateDatabase(name string) error {
	if s.db == nil {
		return mistake.Newf(msg.NOT_DRIVER_DB)
	}

	exist, err := s.ExistDatabase(name)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	sql := jdb.SQLDDL(`	
	CREATE DATABASE $1`, name)
	_, err = jdb.Exec(s.db, sql, name)
	if err != nil {
		return err
	}

	console.Logf(s.name, `Database %s created`, name)

	return nil
}

/**
* DropDatabase
* @param name string
* @return error
**/
func (s *Mysql) DropDatabase(name string) error {
	if s.db == nil {
		return mistake.Newf(msg.NOT_DRIVER_DB)
	}

	exist, err := s.ExistDatabase(name)
	if err != nil {
		return err
	}

	if !exist {
		return nil
	}

	sql := jdb.SQLDDL(`DROP DATABASE $1`, name)
	_, err = jdb.Exec(s.db, sql, name)
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

	s.db, err = s.connectTo(defaultChain)
	if err != nil {
		return nil, err
	}

	params := connection.Params.(*Connection)
	params.Database = connection.Name
	err = s.CreateDatabase(params.Database)
	if err != nil {
		return nil, err
	}

	if s.db != nil {
		err := s.db.Close()
		if err != nil {
			return nil, err
		}
	}

	chain, err := params.Chain()
	if err != nil {
		return nil, err
	}

	s.db, err = s.connectTo(chain)
	if err != nil {
		return nil, err
	}

	s.connected = s.db != nil
	console.Logf(s.name, `Connected to %s:%s`, params.Host, params.Database)

	return s.db, nil
}

/**
* Disconnect
* @return error
**/
func (s *Mysql) Disconnect() error {
	if !s.connected {
		return nil
	}

	if s.db != nil {
		s.db.Close()
	}

	return nil
}
