package jdb

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

const (
	TypeInt      = "int"
	TypeFloat    = "float"
	TypeKey      = "key"
	TypeText     = "text"
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
	RECORDID   = "index"
	STATUS     = "status"
	ACTIVE     = "active"
	ARCHIVED   = "archived"
	CANCELLED  = "cancelled"
	OF_SYSTEM  = "of_system"
	FOR_DELETE = "for_delete"
	TypeData   = map[string]bool{
		TypeInt:      true,
		TypeFloat:    true,
		TypeKey:      true,
		TypeText:     true,
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
		TypeDateTime: true,
		TypeBoolean:  true,
		TypeJson:     true,
		TypeIndex:    true,
		TypeBytes:    true,
		TypeGeometry: true,
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

type DataFunctionTx func(tx *Tx, data et.Json) error

type DataContext func(data et.Json)

type Val struct {
	Value interface{} `json:"value"`
}

/**
* V
* @param value interface{}
* @return *Val
**/
func V(value interface{}) *Val {
	return &Val{
		Value: value,
	}
}

type Column struct {
	Name string `json:"name"`
}

/**
* C
* @param name string
* @return *Column
**/
func C(name string) *Column {
	return &Column{
		Name: name,
	}
}

type Model struct {
	Id           string                  `json:"id"`
	Database     string                  `json:"database"`
	Schema       string                  `json:"schema"`
	Name         string                  `json:"name"`
	Table        string                  `json:"table"`
	Columns      []et.Json               `json:"columns"`
	SourceField  string                  `json:"source_field"`
	RecordField  string                  `json:"record_field"`
	StatusField  string                  `json:"status_field"`
	Details      []et.Json               `json:"details"`
	Masters      []et.Json               `json:"masters"`
	Rollups      []et.Json               `json:"rollups"`   //SQL
	Relations    []et.Json               `json:"relations"` //SQL
	PrimaryKeys  []string                `json:"primary_keys"`
	ForeignKeys  []et.Json               `json:"foreign_keys"`
	Indices      []string                `json:"indices"`
	Required     []string                `json:"required"`
	IsLocked     bool                    `json:"is_locked"`
	Version      int                     `json:"version"`
	db           *Database               `json:"-"`
	details      map[string]*Model       `json:"-"`
	masters      map[string]*Model       `json:"-"`
	isInit       bool                    `json:"-"`
	IsCore       bool                    `json:"is_core"`
	isDebug      bool                    `json:"-"`
	calls        map[string]*DataContext `json:"-"`
	beforeInsert []DataFunctionTx        `json:"-"`
	beforeUpdate []DataFunctionTx        `json:"-"`
	beforeDelete []DataFunctionTx        `json:"-"`
	afterInsert  []DataFunctionTx        `json:"-"`
	afterUpdate  []DataFunctionTx        `json:"-"`
	afterDelete  []DataFunctionTx        `json:"-"`
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
* LoadModel
* @param schema, name string
* @return (*Model, error)
**/
func LoadModel(schema, name string) (*Model, error) {
	id := fmt.Sprintf("%s.%s", schema, name)
	var result Model
	err := loadModel(id, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

/**
* DeleteModel
* @param schema, name string
* @return error
**/
func DeleteModel(schema, name string) error {
	id := fmt.Sprintf("%s.%s", schema, name)
	return deleteModel(id)
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
	bt, err := s.serialize()
	if err != nil {
		return err
	}

	return setModel("model", s.Id, s.Version, bt)
}

/**
* load
* @return error
**/
func (s *Model) load() error {
	return loadModel(s.Id, s)
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
* GetColumnIndex
* @param name string
* @return int
**/
func (s *Model) getColumnIndex(name string) int {
	return slices.IndexFunc(s.Columns, func(item et.Json) bool { return item.String("name") == name })
}

/**
* GetColumn
* @param name string
* @return (et.Json, bool)
**/
func (s *Model) GetColumn(name string) (et.Json, bool) {
	idx := s.getColumnIndex(name)
	if idx == -1 {
		return et.Json{}, false
	}

	result := s.Columns[idx]
	return result, true
}

/**
* SetInit
* @return
**/
func (s *Model) SetInit() {
	s.isInit = true
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
func (s *Model) Debug() {
	s.isDebug = true
}

/**
* Init
* @return error
**/
func (s *Model) Init() error {
	if s.isInit {
		return nil
	}

	err := s.db.init(s)
	if err != nil {
		return err
	}

	for _, m := range s.details {
		m.isDebug = s.isDebug
		err := m.Init()
		if err != nil {
			return err
		}
	}

	return nil
}
