package postgres

import (
	"encoding/json"
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
	"github.com/cgalvisleon/jdb/jdb"
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
* defaultChain
* @return string, error
**/
func (s *Connection) defaultChain() (string, error) {
	return fmt.Sprintf(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, jdb.PostgresDriver, s.Username, s.Password, s.Host, s.Port, "postgres", s.App), nil
}

/**
* chain
* @return string, error
**/
func (s *Connection) chain() (string, error) {
	err := s.validate()
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, jdb.PostgresDriver, s.Username, s.Password, s.Host, s.Port, s.Database, s.App)

	return result, nil
}

/**
* ToJson
* @return et.Json
**/
func (s *Connection) ToJson() et.Json {
	bt, err := json.Marshal(s)
	if err != nil {
		return et.Json{}
	}

	var result et.Json
	err = json.Unmarshal(bt, &result)
	if err != nil {
		return et.Json{}
	}

	return result
}

/**
* Load
* @param params et.Json
* @return error
**/
func (s *Connection) Load(params et.Json) error {
	database := params.Str("database")
	if !utility.ValidStr(database, 0, []string{}) {
		return fmt.Errorf("database is required")
	}

	host := params.Str("host")
	if !utility.ValidStr(host, 0, []string{}) {
		return fmt.Errorf("host is required")
	}

	port := params.Int("port")
	if port == 0 {
		return fmt.Errorf("port is required")
	}

	username := params.Str("username")
	if !utility.ValidStr(username, 0, []string{}) {
		return fmt.Errorf("username is required")
	}

	password := params.Str("password")
	if !utility.ValidStr(password, 0, []string{}) {
		return fmt.Errorf("password is required")
	}

	app := params.Str("app")
	if !utility.ValidStr(app, 0, []string{}) {
		return fmt.Errorf("app is required")
	}

	version := params.Int("version")
	if version == 0 {
		return fmt.Errorf("version is required")
	}

	s.Database = database
	s.Host = host
	s.Port = port
	s.Username = username
	s.Password = password
	s.App = app
	s.Version = version

	return s.validate()
}

/**
* validate
* @return error
**/
func (s *Connection) validate() error {
	if s.Database == "" {
		return fmt.Errorf("database is required")
	}
	if s.Host == "" {
		return fmt.Errorf("host is required")
	}
	if s.Port == 0 {
		return fmt.Errorf("port is required")
	}
	if s.Username == "" {
		return fmt.Errorf("username is required")
	}

	if s.Password == "" {
		return fmt.Errorf("password is required")
	}

	if s.App == "" {
		return fmt.Errorf("app is required")
	}

	return nil
}
