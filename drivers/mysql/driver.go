package mysql

import (
	"database/sql"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	name      string
	params    et.Json
	connStr   string
	db        *sql.DB
	connected bool
	version   int
	nodeId    int
}

func newDriver() jdb.Driver {
	return &Mysql{
		name:      jdb.MysqlDriver,
		params:    et.Json{},
		connected: false,
	}
}

func (s *Mysql) Name() string {
	return s.name
}

func init() {
	jdb.Register(jdb.PostgresDriver, newDriver)
}
