package jdb

import (
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

const (
	TypeInt      = jdb.TypeInt
	TypeFloat    = jdb.TypeFloat
	TypeKey      = jdb.TypeKey
	TypeText     = jdb.TypeText
	TypeMemo     = jdb.TypeMemo
	TypeDateTime = jdb.TypeDateTime
	TypeBoolean  = jdb.TypeBoolean
	TypeJson     = jdb.TypeJson
	TypeIndex    = jdb.TypeIndex
	TypeBytes    = jdb.TypeBytes
	TypeGeometry = jdb.TypeGeometry
	TypeAtribute = jdb.TypeAtribute
	TypeCalc     = jdb.TypeCalc
	TypeDetail   = jdb.TypeDetail
	TypeMaster   = jdb.TypeMaster
	TypeRollup   = jdb.TypeRollup
	TypeRelation = jdb.TypeRelation
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
type Model struct {
	*jdb.Model
}

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

	return &Model{
		Model: result,
	}, nil
}

/**
* DefineColumn
* @param name string, columnType string
* @return error
**/
func (s *Model) DefineColumn(name string, columnType string) error {
	return s.Model.DefineColumn(name, et.Json{
		"type": columnType,
	})
}
