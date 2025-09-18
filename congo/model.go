package jdb

import (
	"encoding/json"
	"fmt"

	"github.com/cgalvisleon/et/envar"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/utility"
)

const (
	SOURCE       = "source"
	KEY          = "id"
	RECORDID     = "index"
	TypeInt      = "int"
	TypeFloat    = "float"
	TypeKey      = "key"
	TypeText     = "text"
	TypeMemo     = "memo"
	TypeDateTime = "datetime"
	TypeBoolean  = "boolean"
	TypeJson     = "json"
	TypeSerial   = "serial"
	TypeBytes    = "bytes"
	TypeGeometry = "geometry"
	TypeAtribute = "atribute"
	TypeCalc     = "calc"
	TypeDetail   = "detail"
	TypeMaster   = "master"
	TypeRollup   = "rollup"
)

var (
	TypeData = map[string]bool{
		TypeInt:      true,
		TypeFloat:    true,
		TypeKey:      true,
		TypeText:     true,
		TypeMemo:     true,
		TypeDateTime: true,
		TypeBoolean:  true,
		TypeJson:     true,
		TypeSerial:   true,
		TypeBytes:    true,
		TypeGeometry: true,
		TypeAtribute: true,
		TypeCalc:     true,
		TypeDetail:   true,
		TypeMaster:   true,
		TypeRollup:   true,
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
		TypeSerial:   true,
		TypeBytes:    true,
		TypeGeometry: true,
	}

	TypeColumnCalc = map[string]bool{
		TypeCalc:   true,
		TypeDetail: true,
		TypeMaster: true,
		TypeRollup: true,
	}
)

type DataFunctionTx func(tx *Tx, data et.Json) error

type Model struct {
	Id           string            `json:"id"`
	Database     string            `json:"database"`
	Schema       string            `json:"schema"`
	Name         string            `json:"name"`
	Table        string            `json:"table"`
	Columns      et.Json           `json:"columns"`
	SourceField  string            `json:"source_field"`
	RecordField  string            `json:"record_field"`
	Details      et.Json           `json:"details"`
	Masters      et.Json           `json:"masters"`
	Rollups      et.Json           `json:"rollups"`
	PrimaryKeys  et.Json           `json:"primary_keys"`
	ForeignKeys  et.Json           `json:"foreign_keys"`
	Indices      []string          `json:"indices"`
	Required     []string          `json:"required"`
	IsLocked     bool              `json:"is_locked"`
	Version      int               `json:"version"`
	db           *Database         `json:"-"`
	details      map[string]*Model `json:"-"`
	masters      map[string]*Model `json:"-"`
	rollups      map[string]*Model `json:"-"`
	isInit       bool              `json:"-"`
	isDebug      bool              `json:"-"`
	beforeInsert []DataFunctionTx  `json:"-"`
	beforeUpdate []DataFunctionTx  `json:"-"`
	beforeDelete []DataFunctionTx  `json:"-"`
	afterInsert  []DataFunctionTx  `json:"-"`
	afterUpdate  []DataFunctionTx  `json:"-"`
	afterDelete  []DataFunctionTx  `json:"-"`
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

	schema := definition.String("schema")
	if !utility.ValidStr(schema, 0, []string{}) {
		return nil, fmt.Errorf(MSG_SCHEMA_REQUIRED)
	}

	name := definition.String("name")
	if !utility.ValidStr(name, 0, []string{}) {
		return nil, fmt.Errorf(MSG_NAME_REQUIRED)
	}

	driver := envar.GetStr("DB_DRIVER", DriverPostgres)
	driver = definition.ValStr(driver, "driver")
	db, err := getDatabase(database, driver)
	if err != nil {
		return nil, err
	}

	result, err := db.getModel(schema, name)
	result.Table = definition.String("table")
	result.Indices = definition.ArrayStr("indices")
	result.Required = definition.ArrayStr("required")
	result.Version = definition.Int("version")

	columns := definition.Json("columns")
	for k := range columns {
		param := columns.Json(k)
		err := result.defineColumn(k, param)
		if err != nil {
			return nil, err
		}
	}

	atribs := definition.Json("atribs")
	for k, v := range atribs {
		err := result.defineAtrib(k, v)
		if err != nil {
			return nil, err
		}
	}

	primaryKeys := definition.ArrayStr("primary_keys")
	result.definePrimaryKeys(primaryKeys...)

	sourceField := definition.String("source_field")
	err = result.defineSourceField(sourceField)
	if err != nil {
		return nil, err
	}

	if err := result.validate(); err != nil {
		return nil, err
	}

	details := definition.Json("details")
	if !details.IsEmpty() {
		err := result.defineDetails(details)
		if err != nil {
			return nil, err
		}
	}

	required := definition.ArrayStr("required")
	result.defineRequired(required...)

	debug := definition.Bool("debug")
	result.isDebug = debug

	return result, nil
}

/**
* ToJson
* @return et.Json
**/
func (s *Model) ToJson() et.Json {
	bt, err := json.Marshal(s)
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
* validate
* @return error
**/
func (s *Model) validate() error {
	if len(s.Columns) != 0 {
		return nil
	}

	err := s.defineSourceField(SOURCE)
	if err != nil {
		return err
	}

	s.defineColumn(KEY, et.Json{
		"type": TypeKey,
	})

	s.definePrimaryKeys(KEY)

	s.defineRecordField(RECORDID)

	return nil
}

/**
* save
* @return error
**/
func (s *Model) save() error {
	return setModel(s.Id, s.ToJson(), s.isDebug)
}

/**
* GetColumn
* @param name string
* @return (et.Json, bool)
**/
func (s *Model) GetColumn(name string) (et.Json, bool) {
	_, ok := s.Columns[name]
	if !ok {
		return et.Json{}, false
	}

	result := s.Columns.Json(name)
	return result, ok
}

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
	if err := s.validate(); err != nil {
		return err
	}

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
