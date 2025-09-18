package jdb

import (
	"database/sql"
	"fmt"
)

func init() {
	Register("postgres", &PostgresDriver{})
}

type PostgresDriver struct {
}

/**
* Connect
* @param db *Database
* @return (*sql.DB, error)
**/
func (s *PostgresDriver) Connect(db *Database) (*sql.DB, error) {
	return nil, nil
}

/**
* Load
* @param model *Model
* @return error
**/
func (s *PostgresDriver) Load(model *Model) error {
	model.Table = fmt.Sprintf("%s.%s", model.Schema, model.Name)

	model.isInit = true

	return nil
}
