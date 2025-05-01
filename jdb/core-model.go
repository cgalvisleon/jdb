package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/timezone"
)

var coreModel *Model

func (s *DB) defineModel() error {
	if err := s.defineSchema(); err != nil {
		return err
	}

	if coreModel != nil {
		return nil
	}

	coreModel = NewModel(coreSchema, "models", 1)
	coreModel.DefineColumn(CREATED_AT, CreatedAtField.TypeData())
	coreModel.DefineColumn(UPDATED_AT, UpdatedAtField.TypeData())
	coreModel.DefineColumn("kind", TypeDataText)
	coreModel.DefineColumn("name", TypeDataText)
	coreModel.DefineColumn("version", TypeDataInt)
	coreModel.DefineColumn("definition", TypeDataBytes)
	coreModel.DefineSystemKeyField()
	coreModel.DefineIndexField()
	coreModel.DefinePrimaryKey("kind", "name")
	coreModel.DefineIndex(true,
		"version",
		SYSID,
		INDEX,
	)
	if err := coreModel.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* getModel
* @param kind, name string
* @return et.Item, error
**/
func (s *DB) getModel(kind, name string) (et.Item, error) {
	return coreModel.
		Where("kind").Eq(kind).
		And("name").Eq(name).
		One()
}

/**
* upsertModel
* @param kind, name string, version int, definition []byte
* @return error
**/
func (s *DB) upsertModel(kind, name string, version int, definition []byte) error {
	now := timezone.Now()
	item, err := coreModel.
		Update(et.Json{
			UPDATED_AT:   now,
			"version":    version,
			"definition": definition,
		}).
		Where("kind").Eq(kind).
		And("name").Eq(name).
		Exec()
	if err != nil {
		return err
	}

	if item.Ok {
		return nil
	}

	_, err = coreModel.
		Insert(et.Json{
			CREATED_AT:   now,
			UPDATED_AT:   now,
			"kind":       kind,
			"name":       name,
			"version":    version,
			"definition": definition,
		}).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

/**
* deleteModel
* @param kind, name string
* @return error
**/
func (s *DB) deleteModel(kind, name string) error {
	_, err := coreModel.
		Delete("kind").Eq(kind).
		And("name").Eq(name).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

/**
* QueryModel
* @param search et.Json
* @return interface{}, error
**/
func (s *DB) QueryModel(search et.Json) (interface{}, error) {
	result, err := coreModel.
		Query(search)
	if err != nil {
		return nil, err
	}

	return result, nil
}
