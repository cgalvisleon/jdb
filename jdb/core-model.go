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
	coreModel.DefineAtribute(CREATED_AT, CreatedAtField.TypeData())
	coreModel.DefineAtribute(UPDATED_AT, UpdatedAtField.TypeData())
	coreModel.DefineAtribute("kind", TypeDataText)
	coreModel.DefineAtribute("name", TypeDataText)
	coreModel.DefineAtribute("version", TypeDataInt)
	coreModel.DefineAtribute("definition", TypeDataBytes)
	coreModel.DefineAtribute(SYSID, SystemKeyField.TypeData())
	coreModel.DefineAtribute(INDEX, IndexField.TypeData())
	coreModel.DefinePrimaryKey(SYSID)
	coreModel.DefineIndex(true,
		"kind",
		"name",
		"version",
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
		Return("version").
		One()
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
		}).Exec()
	if err != nil {
		return err
	}

	return nil
}
