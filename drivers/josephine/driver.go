package josephine

import (
	"database/sql"
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
	_ "modernc.org/sqlite"
)

type Josephine struct {
	arg       et.Json
	connStr   string
	db        *sql.DB
	master    *sql.DB
	connected bool
	version   int
	nodeId    int
}

func NewDriver() jdb.Driver {
	return &Josephine{
		arg:       et.Json{},
		connected: false,
	}
}

func (s *Josephine) Name() string {
	return jdb.Postgres
}

func init() {
	jdb.Register(jdb.Postgres, NewDriver)
}

func (s *Josephine) Disconnect() error {
	return nil
}

func (s *Josephine) SetMain(arg et.Json) error {
	return nil
}

// Database
func (s *Josephine) CreateDatabase(name string) error {
	return nil
}

func (s *Josephine) DropDatabase(name string) error {
	return nil
}

// Core
func (s *Josephine) CreateCore() error {
	return nil
}

// User
func (s *Josephine) GrantPrivileges(username, dbName string) error {
	return nil
}

func (s *Josephine) CreateUser(username, password, confirmation string) error {
	return nil
}

func (s *Josephine) ChangePassword(username, password, confirmation string) error {
	return nil
}

func (s *Josephine) DeleteUser(username string) error {
	return nil
}

// Schema
func (s *Josephine) CreateSchema(name string) error
func (s *Josephine) DropSchema(name string) error

// Model
func (s *Josephine) LoadTable(model *jdb.Model) (bool, error)
func (s *Josephine) CreateModel(model *jdb.Model) error
func (s *Josephine) DropModel(model *jdb.Model) error

// Query
func (s *Josephine) Exec(sql string, arg ...any) error
func (s *Josephine) Query(sql string, arg ...any) (et.Items, error)
func (s *Josephine) One(sql string, arg ...any) (et.Item, error)
func (s *Josephine) Data(source, sql string, arg ...any) (et.Items, error)
func (s *Josephine) Select(ql *jdb.Ql) (et.Items, error)
func (s *Josephine) Count(ql *jdb.Ql) (int, error)
func (s *Josephine) Exists(ql *jdb.Ql) (bool, error)

// Command
func (s *Josephine) Command(command *jdb.Command) (et.Items, error)

// Series
func (s *Josephine) GetSerie(tag string) int64
func (s *Josephine) NextCode(tag, prefix string) string
func (s *Josephine) SetSerie(tag string, val int) int64
func (s *Josephine) CurrentSerie(tag string) int64

// Key Value
func (s *Josephine) SetKey(key string, value []byte) error
func (s *Josephine) GetKey(key string) (et.KeyValue, error)
func (s *Josephine) DeleteKey(key string) error
func (s *Josephine) FindKeys(search string, page, rows int) (et.List, error)

// Function
func (s *Josephine) SetFlow(name string, value []byte) error
func (s *Josephine) GetFlow(id string) (jdb.Flow, error)
func (s *Josephine) DeleteFlow(id string) error
func (s *Josephine) FindFlows(search string, page, rows int) (et.List, error)

// Cache
func (s *Josephine) SetCache(key string, value []byte, duration time.Duration) error
func (s *Josephine) GetCache(key string) (et.KeyValue, error)
func (s *Josephine) DeleteCache(key string) error
func (s *Josephine) CleanCache() error
func (s *Josephine) FindCache(search string, page, rows int) (et.List, error)
