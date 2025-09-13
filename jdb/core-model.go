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
	coreModel.DefineColumn(cf.CreatedAt, TypeDataDateTime)
	coreModel.DefineColumn(cf.UpdatedAt, TypeDataDateTime)
	coreModel.DefineColumn("kind", TypeDataText)
	coreModel.DefineColumn("name", TypeDataText)
	coreModel.DefineColumn("version", TypeDataInt)
	coreModel.DefineColumn("definition", TypeDataBytes)
	coreModel.DefineIndexField()
	coreModel.DefinePrimaryKey("kind", "name")
	coreModel.DefineIndex(true,
		"version",
		cf.SystemId,
		cf.Index,
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
	if coreModel == nil || !coreModel.isInit {
		return et.Item{}, nil
	}

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
	if coreModel == nil || !coreModel.isInit {
		return nil
	}

	now := timezone.Now()
	data := et.Json{
		"kind":       kind,
		"name":       name,
		"version":    version,
		"definition": definition,
	}
	_, err := coreModel.
		Upsert(data).
		BeforeInsert(func(tx *Tx, data et.Json) error {
			data.Set(cf.CreatedAt, now)
			data.Set(cf.UpdatedAt, now)
			return nil
		}).
		BeforeUpdate(func(tx *Tx, data et.Json) error {
			data.Set(cf.UpdatedAt, now)
			return nil
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
	if coreModel == nil || !coreModel.isInit {
		return nil
	}

	_, err := coreModel.
		Delete().
		Where("kind").Eq(kind).
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
	if coreModel == nil || !coreModel.isInit {
		return nil, nil
	}

	result, err := coreModel.
		Query(search)
	if err != nil {
		return nil, err
	}

	return result, nil
}
