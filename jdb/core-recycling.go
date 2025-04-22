package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/timezone"
	"github.com/cgalvisleon/et/utility"
)

var coreRecycling *Model

func (s *DB) defineRecycling() error {
	if err := s.defineSchema(); err != nil {
		return err
	}

	if coreRecycling != nil {
		return nil
	}

	coreRecycling = NewModel(coreSchema, "recycling", 1)
	coreRecycling.DefineAtribute(CREATED_AT, CreatedAtField.TypeData())
	coreRecycling.DefineAtribute(UPDATED_AT, UpdatedAtField.TypeData())
	coreRecycling.DefineAtribute("table", TypeDataText)
	coreRecycling.DefineAtribute(SYSID, SystemKeyField.TypeData())
	coreRecycling.DefineAtribute(INDEX, IndexField.TypeData())
	coreRecycling.DefinePrimaryKey("table", SYSID)
	coreRecycling.DefineIndex(true,
		INDEX,
	)
	if err := coreRecycling.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

func (s *DB) upsertRecycling(table, id, status string) error {
	if status != utility.FOR_DELETE {
		_, err := coreRecycling.
			Delete().
			Where("table").Eq(table).
			And(SYSID).Eq(id).
			Exec()
		if err != nil {
			return err
		}

		return nil
	}

	now := timezone.Now()
	item, err := coreRecycling.
		Update(et.Json{
			UPDATED_AT: now,
			"status":   status,
		}).
		Where("table").Eq(table).
		And(SYSID).Eq(id).
		One()
	if err != nil {
		return err
	}

	if item.Ok {
		return nil
	}

	_, err = coreRecycling.
		Insert(et.Json{
			CREATED_AT: now,
			UPDATED_AT: now,
			"table":    table,
			SYSID:      id,
		}).
		Exec()
	if err != nil {
		return err
	}

	return nil
}
