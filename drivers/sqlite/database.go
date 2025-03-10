package sqlite

import (
	"database/sql"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	_ "modernc.org/sqlite"
)

func (s *SqlLite) connectTo(connStr string) (*sql.DB, error) {
	db, err := sql.Open(s.name, connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (s *SqlLite) connect(params et.Json) error {
	database := params.String("database")
	if database == "" {
		return mistake.New("database is required")
	}

	db, err := s.connectTo(database)
	if err != nil {
		return err
	}

	s.db = db
	s.params = params
	s.connStr = database
	s.connected = s.db != nil
	s.nodeId = params.Int("node_id")
	s.version = 3

	console.Logf(s.name, `Connected to %s:%s`, params.Str("host"), database)

	return nil
}

/**
* Connect
* @param params et.Json
* @return error
**/
func (s *SqlLite) Connect(params et.Json) error {
	err := s.connect(params)
	if err != nil {
		return err
	}

	return nil
}

/**
* CreateDatabase
* @param name string
* @return error
**/
func (s *SqlLite) CreateDatabase(name string) error {
	return mistake.New(MSG_FUNCION_NOT_FOUND)
}

/**
* DropDatabase
* @param name string
* @return error
**/
func (s *SqlLite) DropDatabase(name string) error {
	return mistake.New(MSG_FUNCION_NOT_FOUND)
}

/**
* GrantPrivileges
* @param username string, database string
* @return error
**/
func (s *SqlLite) GrantPrivileges(username, database string) error {
	return mistake.New(MSG_FUNCION_NOT_FOUND)
}

/**
* CreateUser
* @param username string, password string, confirmation string
* @return error
**/
func (s *SqlLite) CreateUser(username, password, confirmation string) error {
	return mistake.New(MSG_FUNCION_NOT_FOUND)
}

/**
* ChangePassword
* @param username string, password string, confirmation string
* @return error
**/
func (s *SqlLite) ChangePassword(username, password, confirmation string) error {
	return mistake.New(MSG_FUNCION_NOT_FOUND)
}

/**
* DeleteUser
* @param username string
* @return error
**/
func (s *SqlLite) DeleteUser(username string) error {
	return mistake.New(MSG_FUNCION_NOT_FOUND)
}
