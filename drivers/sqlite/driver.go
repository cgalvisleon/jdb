package sqlite

import (
	"database/sql"
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
	_ "modernc.org/sqlite"
)

type SqlLite struct {
	params    et.Json
	connStr   string
	db        *sql.DB
	master    *sql.DB
	connected bool
	version   int
	nodeId    int
}

func NewDriver() jdb.Driver {
	return &SqlLite{
		params:    et.Json{},
		connected: false,
	}
}

func (s *SqlLite) Name() string {
	return jdb.Postgres
}

func init() {
	jdb.Register(jdb.Postgres, NewDriver)
}

func (s *SqlLite) Disconnect() error {
	return nil
}

func (s *SqlLite) SetMain(arg et.Json) error {
	return nil
}

// Database
func (s *SqlLite) CreateDatabase(name string) error {
	return nil
}

func (s *SqlLite) DropDatabase(name string) error {
	return nil
}

// Core
func (s *SqlLite) CreateCore() error {
	return nil
}

// User
func (s *SqlLite) GrantPrivileges(username, dbName string) error {
	return nil
}

func (s *SqlLite) CreateUser(username, password, confirmation string) error {
	return nil
}

func (s *SqlLite) ChangePassword(username, password, confirmation string) error {
	return nil
}

func (s *SqlLite) DeleteUser(username string) error {
	return nil
}

// Schema
func (s *SqlLite) CreateSchema(name string) error
func (s *SqlLite) DropSchema(name string) error

// Model
func (s *SqlLite) LoadTable(model *jdb.Model) (bool, error)
func (s *SqlLite) CreateModel(model *jdb.Model) error
func (s *SqlLite) DropModel(model *jdb.Model) error

// Query
func (s *SqlLite) Exec(sql string, arg ...any) error
func (s *SqlLite) Query(sql string, arg ...any) (et.Items, error)
func (s *SqlLite) One(sql string, arg ...any) (et.Item, error)
func (s *SqlLite) Data(source, sql string, arg ...any) (et.Items, error)
func (s *SqlLite) Select(ql *jdb.Ql) (et.Items, error)
func (s *SqlLite) Count(ql *jdb.Ql) (int, error)
func (s *SqlLite) Exists(ql *jdb.Ql) (bool, error)

// Command
func (s *SqlLite) Command(command *jdb.Command) (et.Items, error)

// Series
func (s *SqlLite) GetSerie(tag string) int64
func (s *SqlLite) NextCode(tag, prefix string) string
func (s *SqlLite) SetSerie(tag string, val int) int64
func (s *SqlLite) CurrentSerie(tag string) int64

// Key Value
func (s *SqlLite) SetKey(key string, value []byte) error
func (s *SqlLite) GetKey(key string) (et.KeyValue, error)
func (s *SqlLite) DeleteKey(key string) error
func (s *SqlLite) FindKeys(search string, page, rows int) (et.List, error)

// Function
func (s *SqlLite) SetFlow(name string, value []byte) error
func (s *SqlLite) GetFlow(id string) (jdb.Flow, error)
func (s *SqlLite) DeleteFlow(id string) error
func (s *SqlLite) FindFlows(search string, page, rows int) (et.List, error)

// Cache
func (s *SqlLite) SetCache(key string, value []byte, duration time.Duration) error
func (s *SqlLite) GetCache(key string) (et.KeyValue, error)
func (s *SqlLite) DeleteCache(key string) error
func (s *SqlLite) CleanCache() error
func (s *SqlLite) FindCache(search string, page, rows int) (et.List, error)
