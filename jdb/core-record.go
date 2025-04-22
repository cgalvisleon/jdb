package jdb

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/timezone"
	"github.com/cgalvisleon/et/utility"
)

var coreRecords *Model

func (s *DB) defineRecords() error {
	if err := s.defineSchema(); err != nil {
		return err
	}

	if coreRecords != nil {
		return nil
	}

	coreRecords = NewModel(coreSchema, "records", 1)
	coreRecords.DefineAtribute(CREATED_AT, CreatedAtField.TypeData())
	coreRecords.DefineAtribute(UPDATED_AT, UpdatedAtField.TypeData())
	coreRecords.DefineAtribute("table", TypeDataText)
	coreRecords.DefineAtribute("option", TypeDataText)
	coreRecords.DefineAtribute("sync", TypeDataBool)
	coreRecords.DefineAtribute(SYSID, SystemKeyField.TypeData())
	coreRecords.DefineAtribute(INDEX, IndexField.TypeData())
	coreRecords.DefinePrimaryKey("table", SYSID)
	coreRecords.DefineIndex(true,
		"option",
		"sync",
		SYSID,
		INDEX,
	)
	if err := coreRecords.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

func (s *DB) upsertRecord(table, option, sysid string) error {
	if sysid == "" {
		return mistake.New(MSG_SYSID_REQUIRED)
	}

	current, err := coreRecords.
		Where("table").Eq(table).
		And(SYSID).Eq(sysid).
		One()
	if err != nil {
		return err
	}

	now := timezone.Now()
	if current.Ok {
		_, err := coreRecords.Update(et.Json{
			UPDATED_AT: now,
			"option":   option,
			"sync":     false,
		}).
			Where("table").Eq(table).
			And(SYSID).Eq(sysid).
			One()
		if err != nil {
			return err
		}

		return nil
	}

	_, err = coreRecords.Insert(et.Json{
		CREATED_AT: now,
		UPDATED_AT: now,
		"table":    table,
		"option":   option,
		"sync":     false,
		SYSID:      sysid,
		INDEX:      utility.GenIndex(),
	}).
		One()
	if err != nil {
		return err
	}

	return nil
}

/**
* QueryRecords
* @param query et.Json
* @return interface{}, error
**/
func (s *DB) QueryRecords(query et.Json) (interface{}, error) {
	result, err := coreRecords.
		Query(query)
	if err != nil {
		return nil, err
	}

	return result, nil
}
