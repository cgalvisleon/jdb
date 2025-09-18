package postgres

import (
	"fmt"

	"github.com/cgalvisleon/et/envar"
	jdb "github.com/cgalvisleon/jdb/congo"
)

func init() {
	jdb.Register("postgres", &Postgres{})
}

type Postgres struct {
	name       string        `json:"-"`
	database   *jdb.Database `json:"-"`
	connection *Connection   `json:"-"`
}

func newDriver(database *jdb.Database) jdb.Driver {
	return &Postgres{
		database: database,
		name:     database.Name,
		connection: &Connection{
			Database: envar.GetStr("DB_NAME", "jdb"),
			Host:     envar.GetStr("DB_HOST", "localhost"),
			Port:     envar.GetInt("DB_PORT", 5432),
			Username: envar.GetStr("DB_USER", "admin"),
			Password: envar.GetStr("DB_PASSWORD", "admin"),
			App:      envar.GetStr("APP_NAME", "jdb"),
			Version:  envar.GetInt("DB_VERSION", 13),
		},
	}
}

/**
* Load
* @param model *Model
* @return error
**/
func (s *Postgres) Load(model *jdb.Model) error {
	model.Table = fmt.Sprintf("%s.%s", model.Schema, model.Name)
	model.SetInit()

	return nil
}
