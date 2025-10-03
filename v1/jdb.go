package jdb

import (
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

const (
	TypeInt        = jdb.TypeInt
	TypeFloat      = jdb.TypeFloat
	TypeKey        = jdb.TypeKey
	TypeText       = jdb.TypeText
	TypeMemo       = jdb.TypeMemo
	TypeDateTime   = jdb.TypeDateTime
	TypeBoolean    = jdb.TypeBoolean
	TypeJson       = jdb.TypeJson
	TypeIndex      = jdb.TypeIndex
	TypeBytes      = jdb.TypeBytes
	TypeGeometry   = jdb.TypeGeometry
	TypeAtribute   = jdb.TypeAtribute
	TypeCalc       = jdb.TypeCalc
	TypeVm         = jdb.TypeVm
	TypeDetail     = jdb.TypeDetail
	TypeMaster     = jdb.TypeMaster
	TypeRollup     = jdb.TypeRollup
	TypeRelation   = jdb.TypeRelation
	DriverPostgres = jdb.DriverPostgres
	DriverMysql    = jdb.DriverMysql
	DriverSqlite   = jdb.DriverSqlite
	DriverMssql    = jdb.DriverMssql
	DriverOracle   = jdb.DriverOracle
)

var (
	SOURCE         = jdb.SOURCE
	KEY            = jdb.KEY
	RECORDID       = jdb.RECORDID
	STATUS         = jdb.STATUS
	ACTIVE         = jdb.ACTIVE
	ARCHIVED       = jdb.ARCHIVED
	CANCELLED      = jdb.CANCELLED
	OF_SYSTEM      = jdb.OF_SYSTEM
	FOR_DELETE     = jdb.FOR_DELETE
	PENDING        = jdb.PENDING
	APPROVED       = jdb.APPROVED
	REJECTED       = jdb.REJECTED
	CREATED_AT     = jdb.CREATED_AT
	UPDATED_AT     = jdb.UPDATED_AT
	TEAM_ID        = jdb.TEAM_ID
	TypeData       = jdb.TypeData
	TypeColumn     = jdb.TypeColumn
	TypeAtrib      = jdb.TypeAtrib
	TypeColumnCalc = jdb.TypeColumnCalc
	ErrNotInserted = jdb.ErrNotInserted
	ErrNotUpdated  = jdb.ErrNotUpdated
	ErrNotFound    = jdb.ErrNotFound
	ErrNotUpserted = jdb.ErrNotUpserted
	ErrDuplicate   = jdb.ErrDuplicate
)

type DB = jdb.Database
type Model = jdb.Model
type Tx = jdb.Tx
type Condition = jdb.Condition
type Ql = jdb.Ql
type Cmd = jdb.Cmd

/**
* NewModel
* @param db *jdb.Database, schema, name string, version int
* @return (*Model, error)
**/
func NewModel(db *jdb.Database, schema, name string, version int) (*Model, error) {
	result, err := db.Define(et.Json{
		"schema":  schema,
		"name":    name,
		"version": version,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

/**
* ConnectTo
* @param name, driver string, userCore bool, params et.Json
* @return (*jdb.Database, error)
**/
func ConnectTo(name, driver string, userCore bool, params et.Json) (*jdb.Database, error) {
	return jdb.ConnectTo(name, driver, userCore, params)
}

/**
* LoadTo
* @param name string
* @return (*jdb.Database, error)
**/
func LoadTo(name string) (*jdb.Database, error) {
	return jdb.LoadTo(name)
}

/**
* Load
* @return (*jdb.Database, error)
**/
func Load() (*jdb.Database, error) {
	return jdb.Load()
}

/**
* Define
* @param definition et.Json
* @return (*jdb.Model, error)
**/
func Define(definition et.Json) (*jdb.Model, error) {
	return jdb.Define(definition)
}

/**
* Select
* @param query et.Json
* @return (*jdb.Ql, error)
**/
func Select(query et.Json) (*jdb.Ql, error) {
	return jdb.Select(query)
}

/**
* From
* @param model *Model
* @return *Ql
**/
func From(model *Model) *jdb.Ql {
	return jdb.From(model)
}

/**
* Insert
* @param param et.Json
* @return (*jdb.Cmd, error)
**/
func Insert(param et.Json) (*jdb.Cmd, error) {
	return jdb.Insert(param)
}

/**
* Update
* @param param et.Json
* @return (*jdb.Cmd, error)
**/
func Update(param et.Json) (*jdb.Cmd, error) {
	return jdb.Update(param)
}

/**
* Delete
* @param param et.Json
* @return (*jdb.Cmd, error)
**/
func Delete(param et.Json) (*jdb.Cmd, error) {
	return jdb.Delete(param)
}

/**
* Upsert
* @param param et.Json
* @return (*jdb.Cmd, error)
**/
func Upsert(param et.Json) (*jdb.Cmd, error) {
	return jdb.Upsert(param)
}
