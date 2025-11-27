package jdb

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/reg"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
)

const (
	TypeInt      = "int"
	TypeFloat    = "float"
	TypeKey      = "key"
	TypeText     = "text"
	TypeMemo     = "memo"
	TypeDateTime = "datetime"
	TypeBoolean  = "boolean"
	TypeJson     = "json"
	TypeIndex    = "index"
	TypeBytes    = "bytes"
	TypeGeometry = "geometry"
	TypeAtribute = "atribute"
	TypeCalc     = "calc"
	TypeDetail   = "detail"
	TypeMaster   = "master"
	TypeRollup   = "rollup"
	TypeRelation = "relation"
)

var (
	SOURCE     = "source"
	KEY        = "id"
	RECORDID   = "idx"
	STATUS     = "status"
	ACTIVE     = "active"
	ARCHIVED   = "archived"
	CANCELLED  = "cancelled"
	OF_SYSTEM  = "of_system"
	FOR_DELETE = "for_delete"
	PENDING    = "pending"
	APPROVED   = "approved"
	REJECTED   = "rejected"
	CREATED_AT = "created_at"
	UPDATED_AT = "updated_at"
	TENANT_ID  = "tenant_id"
	PROJECT_ID = "project_id"

	TypeData = map[string]bool{
		TypeInt:      true,
		TypeFloat:    true,
		TypeKey:      true,
		TypeText:     true,
		TypeMemo:     true,
		TypeDateTime: true,
		TypeBoolean:  true,
		TypeJson:     true,
		TypeIndex:    true,
		TypeBytes:    true,
		TypeGeometry: true,
		TypeAtribute: true,
		TypeCalc:     true,
		TypeDetail:   true,
		TypeMaster:   true,
		TypeRollup:   true,
		TypeRelation: true,
	}

	TypeColumn = map[string]bool{
		TypeInt:      true,
		TypeFloat:    true,
		TypeKey:      true,
		TypeText:     true,
		TypeMemo:     true,
		TypeDateTime: true,
		TypeBoolean:  true,
		TypeJson:     true,
		TypeIndex:    true,
		TypeBytes:    true,
		TypeGeometry: true,
	}

	TypeAtrib = map[string]bool{
		TypeAtribute: true,
	}

	TypeColumnCalc = map[string]bool{
		TypeCalc:     true,
		TypeDetail:   true,
		TypeMaster:   true,
		TypeRollup:   true,
		TypeRelation: true,
	}

	ErrNotInserted = fmt.Errorf("record not inserted")
	ErrNotUpdated  = fmt.Errorf("record not updated")
	ErrNotFound    = fmt.Errorf("record not found")
	ErrNotUpserted = fmt.Errorf("record not inserted or updated")
	ErrDuplicate   = fmt.Errorf("record duplicate")
)

type DataFunctionTx func(tx *Tx, old, new et.Json) error
type DataContext func(data et.Json)

type Model struct {
	Database      string                 `json:"database"`
	Schema        string                 `json:"schema"`
	Name          string                 `json:"name"`
	Table         string                 `json:"table"`
	Columns       []*Column              `json:"columns"`
	Hidden        []string               `json:"hidden"`
	SourceField   string                 `json:"source_field"`
	RecordField   string                 `json:"record_field"`
	StatusField   string                 `json:"status_field"`
	Details       map[string]*Detail     `json:"details"`
	Rollups       map[string]*Detail     `json:"rollups"`
	Relations     map[string]*Detail     `json:"relations"`
	Calcs         map[string]DataContext `json:"-"`
	UniqueIndexes []string               `json:"unique_indexes"`
	PrimaryKeys   []string               `json:"primary_keys"`
	ForeignKeys   []et.Json              `json:"foreign_keys"`
	Indexes       []string               `json:"indexes"`
	Required      []string               `json:"required"`
	IsLocked      bool                   `json:"is_locked"`
	Version       int                    `json:"version"`
	Current       int                    `json:"current"`
	IsDebug       bool                   `json:"-"`
	db            *DB                    `json:"-"`
	isInit        bool                   `json:"-"`
	isCore        bool                   `json:"-"`
	beforeInserts []DataFunctionTx       `json:"-"`
	beforeUpdates []DataFunctionTx       `json:"-"`
	beforeDeletes []DataFunctionTx       `json:"-"`
	afterInserts  []DataFunctionTx       `json:"-"`
	afterUpdates  []DataFunctionTx       `json:"-"`
	afterDeletes  []DataFunctionTx       `json:"-"`
}

/**
* Define
* @param definition et.Json
* @return (*Model, error)
**/
func Define(definition et.Json) (*Model, error) {
	database := definition.String("database")
	if !utility.ValidStr(database, 0, []string{}) {
		return nil, fmt.Errorf(MSG_DATABASE_REQUIRED)
	}

	driver := envar.GetStr("DB_DRIVER", DriverPostgres)
	driver = definition.ValStr(driver, "driver")
	connection := definition.Json("connection")
	userCore := definition.Bool("user_core")
	db, err := getDatabase(database, driver, userCore, connection)
	if err != nil {
		return nil, err
	}

	result, err := db.Define(definition)
	if err != nil {
		return nil, err
	}

	return result, nil
}

/**
* DeleteModel
* @param name string
* @return error
**/
func DeleteModel(name string) error {
	return deleteModel(name)
}

/**
* Serialize
* @return ([]byte, error)
**/
func (s *Model) serialize() ([]byte, error) {
	bt, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return bt, nil
}

/**
* save
* @return error
**/
func (s *Model) save() error {
	if s.isCore {
		return nil
	}

	bt, err := s.serialize()
	if err != nil {
		return err
	}

	return setModel(s.Name, s.Version, bt)
}

/**
* prepare
* @return error
**/
func (s *Model) prepare() error {
	if len(s.Columns) != 0 {
		return nil
	}

	s.DefineModel()

	return nil
}

/**
* GetColumnIndex
* @param name string
* @return int
**/
func (s *Model) getColumnIndex(name string) int {
	return slices.IndexFunc(s.Columns, func(item *Column) bool { return item.Name == name })
}

/**
* UseAtribs
* @return bool
**/
func (s *Model) UseAtribs() bool {
	return s.SourceField != "" && !s.IsLocked
}

/**
* ToJson
* @return et.Json
**/
func (s *Model) ToJson() et.Json {
	bt, err := s.serialize()
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
* GetColumn
* @param name string
* @return (*Column, bool)
**/
func (s *Model) GetColumn(name string) (*Column, bool) {
	spPeriod := strings.Split(name, ".")
	if len(spPeriod) > 1 {
		spColon := strings.Split(spPeriod[1], ":")
		if len(spColon) > 1 {
			return nil, false
		}
		name = spPeriod[1]
	}
	idx := s.getColumnIndex(name)
	if idx == -1 {
		return nil, false
	}

	result := s.Columns[idx]
	return result, true
}

/**
* Init
* @return error
**/
func (s *Model) Init() error {
	if s.isInit {
		return nil
	}

	s.Current = versionModel(s.Name, s.Version)
	err := s.db.init(s)
	if err != nil {
		return err
	}

	for _, m := range s.Details {
		m.From.IsDebug = s.IsDebug
		err := m.From.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

/**
* Lock
* @return
**/
func (s *Model) Lock() {
	s.IsLocked = true
	s.save()
}

/**
* Unlock
* @return
**/
func (s *Model) Unlock() {
	s.IsLocked = false
	s.save()
}

/**
* Debug
* @return
**/
func (s *Model) Debug() *Model {
	s.IsDebug = true
	return s
}

/**
* GetId
* @param id string
* @return string
**/
func (s *Model) GetId(id string) string {
	if !map[string]bool{"": true, "*": true, "new": true}[id] {
		return id
	}

	return reg.GenULID(s.Name)
}

/**
* GetRecordById
* @param id string
* @return (et.Item, error)
**/
func (s *Model) GetRecordById(id string) (et.Item, error) {
	result, err := s.
		Where(Eq(RECORDID, id)).
		One()
	if err != nil {
		return et.Item{}, err
	}

	if !result.Ok {
		return et.Item{}, fmt.Errorf(MSG_RECORD_NOT_FOUND, id)
	}

	return result, nil
}

/**
* GetKey
* @param data et.Json
* @return string
**/
func (s *Model) GetKey(data et.Json) string {
	result := ""
	for _, col := range s.PrimaryKeys {
		value := data.Str(col)
		result = strs.Append(result, value, ":")
	}

	return result
}
