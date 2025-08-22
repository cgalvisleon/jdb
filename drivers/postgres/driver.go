package postgres

import (
	"database/sql"
	"errors"

	"github.com/cgalvisleon/et/config"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
	_ "github.com/lib/pq"
)

type Connection struct {
	Database string `json:"database"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	App      string `json:"app"`
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

	result := strs.Format(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, jdb.PostgresDriver, s.Username, s.Password, s.Host, s.Port, s.Database, s.App)

	return result, nil
}

/**
* defaultChain
* @return string, error
**/
func (s *Connection) defaultChain() (string, error) {
	return strs.Format(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, jdb.PostgresDriver, s.Username, s.Password, s.Host, s.Port, "postgres", s.App), nil
}

/**
* ToJson
* @return et.Json
**/
func (s *Connection) ToJson() et.Json {
	return et.Json{
		"database": s.Database,
		"host":     s.Host,
		"port":     s.Port,
		"username": s.Username,
		"password": s.Password,
		"app":      s.App,
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
	if s.Host == "" {
		return errors.New("host is required")
	}
	if s.Port == 0 {
		return errors.New("port is required")
	}
	if s.Username == "" {
		return errors.New("username is required")
	}

	if s.Password == "" {
		return errors.New("password is required")
	}

	if s.App == "" {
		return errors.New("app is required")
	}

	return nil
}

type Postgres struct {
	db         *sql.DB
	name       string
	version    int
	connected  bool
	connection Connection
}

func newDriver() jdb.Driver {
	return &Postgres{
		name: jdb.PostgresDriver,
		connection: Connection{
			Database: config.String("DB_NAME", "jdb"),
			Host:     config.String("DB_HOST", "localhost"),
			Port:     config.Int("DB_PORT", 5432),
			Username: config.String("DB_USER", "admin"),
			Password: config.String("DB_PASSWORD", "admin"),
			App:      config.App.Name,
			Version:  config.Int("DB_VERSION", 13),
		},
	}
}

func (s *Postgres) Name() string {
	return s.name
}

func init() {
	jdb.Register(jdb.PostgresDriver, newDriver, jdb.ConnectParams{
		Id:       config.String("DB_ID", "jdb"),
		Driver:   jdb.PostgresDriver,
		Name:     config.String("DB_NAME", "jdb"),
		UserCore: true,
		NodeId:   config.Int("NODE_ID", 0),
		Debug:    config.Bool("DEBUG", false),
		Params: &Connection{
			Database: config.String("DB_NAME", "jdb"),
			Host:     config.String("DB_HOST", "localhost"),
			Port:     config.Int("DB_PORT", 5432),
			Username: config.String("DB_USER", "admin"),
			Password: config.String("DB_PASSWORD", "admin"),
			App:      config.App.Name,
			Version:  config.Int("DB_VERSION", 13),
		},
	})
}
