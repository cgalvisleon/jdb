package sqlite

import (
	"database/sql"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/jdb/jdb"
	_ "modernc.org/sqlite"
)

func (s *SqlLite) connectTo(connStr string) (*sql.DB, error) {
	db, err := sql.Open(jdb.SqlLite, connStr)
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

	database := params.String("database")
	console.Logf(jdb.Postgres, `Connected to %s`, database)

	return nil
}
