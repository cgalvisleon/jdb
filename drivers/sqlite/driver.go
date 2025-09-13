package sqlite

import (
	"errors"

	"github.com/cgalvisleon/et/config"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
	_ "modernc.org/sqlite"
)

type Connection struct {
	Database string `json:"database"`
	Version  int    `json:"version"`
}

/**
* Chain
* @return string, error
**/
func (s *Connection) Chain() (string, error) {
	err := s.Validate()
	if err != nil {
		return "", err
	}

	result := s.Database

	return result, nil
}

/**
* ToJson
* @return et.Json
**/
func (s *Connection) ToJson() et.Json {
	return et.Json{
		"database": s.Database,
		"version":  s.Version,
	}
}

/**
* Validate
* @return error
**/
func (s *Connection) Validate() error {
	if s.Database == "" {
		return errors.New("database is required")
	}

	return nil
}

type SqlLite struct {
	jdb        *jdb.DB
	name       string
	version    int
	connected  bool
	connection Connection
}

func newDriver(db *jdb.DB) jdb.Driver {
	return &SqlLite{
		jdb:       db,
		name:      jdb.SqliteDriver,
		connected: false,
		connection: Connection{
			Database: config.String("DB_NAME", "jdb"),
		},
	}
}

func (s *SqlLite) Name() string {
	return s.name
}

func init() {
	jdb.Register(jdb.SqliteDriver, newDriver, jdb.ConnectParams{
		Id:       config.String("DB_ID", "jdb"),
		Driver:   jdb.SqliteDriver,
		Name:     config.String("DB_NAME", "jdb"),
		UserCore: true,
		NodeId:   config.Int("NODE_ID", 1),
		Debug:    config.Bool("DEBUG", false),
		Params: &Connection{
			Database: config.String("DB_NAME", "jdb"),
		},
	})
}
