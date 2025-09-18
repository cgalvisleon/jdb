package postgres

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
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
	return fmt.Sprintf(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, driver, s.Username, s.Password, s.Host, s.Port, "postgres", s.App), nil
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

	result := fmt.Sprintf(`%s://%s:%s@%s:%d/%s?sslmode=disable&application_name=%s`, driver, s.Username, s.Password, s.Host, s.Port, s.Database, s.App)

	return result, nil
}

/**
* load
* @param params et.Json
* @return error
**/
func (s *Connection) load(params et.Json) error {
	database := params.Str("database")
	if utility.ValidStr(database, 0, []string{}) {
		s.Database = database
	}

	host := params.Str("host")
	if utility.ValidStr(host, 0, []string{}) {
		s.Host = host
	}

	port := params.Int("port")
	if port != 0 {
		s.Port = port
	}

	username := params.Str("username")
	if utility.ValidStr(username, 0, []string{}) {
		s.Username = username
	}

	password := params.Str("password")
	if utility.ValidStr(password, 0, []string{}) {
		s.Password = password
	}

	app := params.Str("app")
	if utility.ValidStr(app, 0, []string{}) {
		s.App = app
	}

	version := params.Int("version")
	if version != 0 {
		s.Version = version
	}

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
