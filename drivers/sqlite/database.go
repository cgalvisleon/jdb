package sqlite

import (
	"database/sql"
	"strings"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/jdb/jdb"
)

func (s *SqlLite) connectTo(database string) (*sql.DB, error) {
	if !strings.HasSuffix(database, ".db") {
		database = database + ".db"
	}

	db, err := sql.Open(s.name, database)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

/**
* Connect
* @param connection jdb.ConnectParams
* @return error
**/
func (s *SqlLite) Connect(connection jdb.ConnectParams) (*sql.DB, error) {
	database := connection.Params.(*Connection).Database
	if database == "" {
		return nil, mistake.New("database is required")
	}

	db, err := s.connectTo(database)
	if err != nil {
		return nil, err
	}

	s.connected = db != nil
	console.Logf(s.name, `Connected to %s`, database)

	return db, nil
}
