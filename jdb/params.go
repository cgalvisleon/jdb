package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
)

type ConnectParams struct {
	Driver   string  `json:"driver"`
	Name     string  `json:"name"`
	Params   et.Json `json:"params"`
	UserCore bool    `json:"user_core"`
}

/**
* Validate
* @return error
**/
func (s *ConnectParams) Validate() error {
	if conn == nil {
		return mistake.New(MSG_JDB_NOT_DEFINED)
	}

	if s.Driver == "" {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	if s.Name == "" {
		return mistake.New(MSG_DATABASE_NOT_DEFINED)
	}

	if _, ok := conn.Drivers[s.Driver]; !ok {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return nil
}
