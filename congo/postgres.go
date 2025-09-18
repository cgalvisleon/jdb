package jdb

import (
	"fmt"
)

func init() {
	Register("postgres", &PostgresDriver{})
}

type PostgresDriver struct {
}

func (s *PostgresDriver) Load(model *Model) (string, error) {
	model.Table = fmt.Sprintf("%s.%s", model.Schema, model.Name)

	return model.ToJson().ToString(), nil
}
