package jdb

import (
	"time"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
)

func Name(name string) string {
	result := strs.ReplaceAll(name, []string{" "}, "_")
	return strs.Lowcase(result)
}

var JDBS []*DB = []*DB{}

type Mode int

const (
	Origin Mode = iota
	Local
)

type DB struct {
	CreatedAt   time.Time          `json:"created_date"`
	UpdateAt    time.Time          `json:"update_date"`
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Schemas     map[string]*Schema `json:"schemas"`
	UseCore     bool               `json:"use_core"`
	NodeId      int64              `json:"node_id"`
	Mode        Mode               `json:"mode"`
	Origin      string             `json:"origin"`
	driver      Driver             `json:"-"`
}

/**
* NewDatabase
* @param name, driver string, id int
* @return *DB
**/
func NewDatabase(name, driver string, id int64) (*DB, error) {
	if driver == "" {
		return nil, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	if _, ok := Jdb.Drivers[driver]; !ok {
		return nil, mistake.Newf(MSG_DRIVER_NOT_FOUND, driver)
	}

	now := time.Now()
	result := &DB{
		CreatedAt:   now,
		UpdateAt:    now,
		Id:          utility.RecordId("db", ""),
		Name:        Name(name),
		Description: "",
		Schemas:     map[string]*Schema{},
		NodeId:      id,
		driver:      Jdb.Drivers[driver](),
	}

	JDBS = append(JDBS, result)

	return result, nil
}

/**
* Describe
* @return et.Json
**/
func (s *DB) Describe() et.Json {
	var schemas = make([]et.Json, 0)
	for _, schema := range s.Schemas {
		schemas = append(schemas, schema.Describe())
	}

	result := et.Json{
		"created_date": s.CreatedAt,
		"update_date":  s.UpdateAt,
		"id":           s.Id,
		"name":         s.Name,
		"description":  s.Description,
		"schemas":      schemas,
	}

	return result
}

/**
* GenId
* @param tag string
* @return string
**/
func (s *DB) GenId(tag string) string {
	return utility.Snowflake(s.NodeId, tag)
}

/**
* Conected
* @return bool
**/
func (s *DB) Conected(params et.Json) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Connect(params)
}

/**
* Disconected
* @return error
**/
func (s *DB) Disconected() error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Disconnect()
}

/**
* SetMain
* @param params et.Json
* @return error
**/
func (s *DB) SetMain(params et.Json) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.SetMain(params)
}

/**
* GrantPrivileges
* @return error
**/
func (s *DB) GrantPrivileges(username, database string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.GrantPrivileges(username, database)
}

/**
* CreateUser
* @return error
**/
func (s *DB) CreateUser(username, password, confirmation string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.CreateUser(username, password, confirmation)
}

/**
* ChangePassword
* @return error
**/
func (s *DB) ChangePassword(username, password, confirmation string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.ChangePassword(username, password, confirmation)
}

/**
* DeleteUser
* @return error
**/
func (s *DB) DeleteUser(username string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.DeleteUser(username)
}

/**
* CreateSchema
* @param name string
* @return error
**/
func (s *DB) CreateSchema(name string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.CreateSchema(name)
}

/**
* DropSchema
* @param name string
* @return error
**/
func (s *DB) DropSchema(name string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.DropSchema(name)
}

/**
* CreateCore
* @return error
**/
func (s *DB) CreateCore() error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.CreateCore()
}

/**
* LoadTable
* @param model *Model
* @return error
**/
func (s *DB) LoadTable(model *Model) (bool, error) {
	if s.driver == nil {
		return false, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.LoadTable(model)
}

/**
* CreateModel
* @param model *Model
* @return error
**/
func (s *DB) CreateModel(model *Model) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.CreateModel(model)
}

/**
* DropModel
* @param model *Model
**/
func (s *DB) DropModel(model *Model) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.DropModel(model)
}

/**
* SaveModel
* @param model *Model
* @return error
**/
func (s *DB) SaveModel(model *Model) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.SaveModel(model)
}

/**
* Exec
* @param sql string
* @param params ...any
* @return error
**/
func (s *DB) Exec(sql string, params ...any) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Exec(sql, params...)
}

/**
* Query
* @param sql string
* @param params ...any
* @return et.Items, error
**/
func (s *DB) Query(sql string, params ...any) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Query(sql, params...)
}

/**
* One
* @param sql string
* @param params ...any
* @return et.Item, error
**/
func (s *DB) One(sql string, params ...any) (et.Item, error) {
	if s.driver == nil {
		return et.Item{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.One(sql, params...)
}

/**
* Data
* @param source string
* @param sql string
* @param params ...any
* @return et.Item, error
**/
func (s *DB) Data(source, sql string, params ...any) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Data(source, sql, params...)
}

/**
* Select
* @param ql *Ql
* @return et.Items, error
**/
func (s *DB) Select(ql *Ql) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Select(ql)
}

/**
* Exists
* @param ql *Ql
* @return bool, error
**/
func (s *DB) Exists(ql *Ql) (bool, error) {
	if s.driver == nil {
		return false, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Exists(ql)
}

/**
* Count
* @param ql *Ql
* @return int, error
**/
func (s *DB) Count(ql *Ql) (int, error) {
	if s.driver == nil {
		return 0, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Count(ql)
}

/**
* Command
* @param command *Command
* @return et.Item, error
**/
func (s *DB) Command(command *Command) (et.Items, error) {
	if s.driver == nil {
		return et.Items{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.Command(command)
}

/**
* GetSerie
* @return int64
**/
func (s *DB) GetSerie(tag string) int64 {
	if s.driver == nil {
		return 0
	}

	return s.driver.GetSerie(tag)
}

/**
* NextCode
* @param tag string
* @param format string "%08v" "PREFIX-%08v-SUFFIX"
* @return string
**/
func (s *DB) NextCode(tag, format string) string {
	if s.driver == nil {
		return ""
	}

	return s.driver.NextCode(tag, format)
}

/**
* SetSerie
* @return int64
**/
func (s *DB) SetSerie(tag string, val int) int64 {
	if s.driver == nil {
		return 0
	}

	return s.driver.SetSerie(tag, val)
}

/**
* CurrentSerie
* @param tag string
* @return int64
**/
func (s *DB) CurrentSerie(tag string) int64 {
	if s.driver == nil {
		return 0
	}

	return s.driver.CurrentSerie(tag)
}

/**
* SetKey
* @param key string
* @param value value []byte
* @return error
**/
func (s *DB) SetKey(key string, value []byte) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.SetKey(key, value)
}

/**
* GetKey
* @param key string
* @return et.KeyValue, error
**/
func (s *DB) GetKey(key string) (et.KeyValue, error) {
	if s.driver == nil {
		return et.KeyValue{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.GetKey(key)
}

/**
* DeleteKey
* @param key string
* @return error
**/
func (s *DB) DeleteKey(key string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.DeleteKey(key)
}

/**
* FindKeys
* @param search string
* @param page int
* @param rows int
* @return et.List, error
**/
func (s *DB) FindKeys(search string, page, rows int) (et.List, error) {
	if s.driver == nil {
		return et.List{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.FindKeys(search, page, rows)
}

/**
* SetFlow
* @param name string, value []byte
* @return error
**/
func (s *DB) SetFlow(name string, value []byte) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.SetFlow(name, value)
}

/**
* GetFlow
* @param name string
* @return Flow, error
**/
func (s *DB) GetFlow(name string) (Flow, error) {
	if s.driver == nil {
		return Flow{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.GetFlow(name)
}

/**
* DeleteFlow
* @param name string
* @return error
**/
func (s *DB) DeleteFlow(name string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.DeleteFlow(name)
}

/**
* FindFlows
* @param search string, page int, rows int
* @return et.List, error
**/
func (s *DB) FindFlows(search string, page, rows int) (et.List, error) {
	if s.driver == nil {
		return et.List{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.FindFlows(search, page, rows)
}

/**
* SetCache
* @param key string, value []byte, duration time.Duration
* @return error
**/
func (s *DB) SetCache(key string, value []byte, duration time.Duration) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.SetCache(key, value, duration)
}

/**
* GetCache
* @param key string
* @return et.KeyValue, error
**/
func (s *DB) GetCache(key string) (et.KeyValue, error) {
	if s.driver == nil {
		return et.KeyValue{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.GetCache(key)
}

/**
* DeleteCache
* @param key string
* @return error
**/
func (s *DB) DeleteCache(key string) error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.DeleteCache(key)
}

/**
* CleanCache
* @return error
**/
func (s *DB) CleanCache() error {
	if s.driver == nil {
		return mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.CleanCache()
}

/**
* FindCache
* @param search string, page int, rows int
* @return et.List, error
**/
func (s *DB) FindCache(search string, page, rows int) (et.List, error) {
	if s.driver == nil {
		return et.List{}, mistake.New(MSG_DRIVER_NOT_DEFINED)
	}

	return s.driver.FindCache(search, page, rows)
}
