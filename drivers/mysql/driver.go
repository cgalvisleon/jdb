package mysql

import (
	"errors"

	"github.com/cgalvisleon/et/config"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/jdb/jdb"
	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	Database string `json:"database"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
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

	result := strs.Format(`%s:%s@tcp(%s:%d)/%s?parseTime=true`, s.Username, s.Password, s.Host, s.Port, s.Database)

	return result, nil
}

/**
* defaultChain
* @return string, error
**/
func (s *Connection) defaultChain() (string, error) {
	return strs.Format(`%s:%s@tcp(%s:%d)/%s?parseTime=true`, s.Username, s.Password, s.Host, s.Port, "mysql"), nil
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

	return nil
}

type Mysql struct {
	jdb        *jdb.DB
	name       string
	version    int
	connected  bool
	connection Connection
}

func newDriver(db *jdb.DB) jdb.Driver {
	return &Mysql{
		jdb:       db,
		name:      jdb.MysqlDriver,
		connected: false,
		connection: Connection{
			Database: config.String("DB_NAME", "jdb"),
			Host:     config.String("DB_HOST", "localhost"),
			Port:     config.Int("DB_PORT", 3306),
			Username: config.String("DB_USER", "admin"),
			Password: config.String("DB_PASSWORD", "admin"),
			Version:  config.Int("DB_VERSION", 8),
		},
	}
}

func (s *Mysql) Name() string {
	return s.name
}

func init() {
	jdb.Register(jdb.PostgresDriver, newDriver, jdb.ConnectParams{
		Id:       config.String("DB_ID", "jdb"),
		Driver:   jdb.OracleDriver,
		Name:     config.String("DB_NAME", "jdb"),
		UserCore: true,
		NodeId:   config.Int("NODE_ID", 1),
		Debug:    config.Bool("DEBUG", false),
		Params: &Connection{
			Database: config.String("DB_NAME", "jdb"),
			Host:     config.String("DB_HOST", "localhost"),
			Port:     config.Int("DB_PORT", 3306),
			Username: config.String("DB_USER", "admin"),
			Password: config.String("DB_PASSWORD", "admin"),
			Version:  config.Int("DB_VERSION", 8),
		},
	})
}
