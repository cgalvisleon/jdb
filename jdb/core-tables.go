package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
)

var coreTables *Model

func (s *DB) defineTables() error {
	if err := s.defineSchema(); err != nil {
		return err
	}

	if coreTables != nil {
		return nil
	}

	coreTables = NewModel(coreSchema, "tables", 1)
	coreTables.DefineColumn("schema_name", TypeDataText)
	coreTables.DefineColumn("table_name", TypeDataText)
	coreTables.DefineColumn("count", TypeDataInt)
	coreTables.DefineIndexField()
	coreTables.DefinePrimaryKey("schema_name", "table_name")
	coreTables.DefineIndex(true,
		"schema_name",
		"table_name",
	)
	if err := coreTables.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

func (s *DB) upsertTable(tx *Tx, schema, name string, inc int) error {
	if coreTables == nil || !coreTables.isInit {
		return nil
	}

	_, err := coreTables.
		Upsert(et.Json{
			"schema_name": schema,
			"table_name":  name,
		}).
		BeforeInsert(func(tx *Tx, data et.Json) error {
			data.Set("count", 1)
			return nil
		}).
		BeforeUpdate(func(tx *Tx, data et.Json) error {
			data.Set("count", fmt.Sprintf(`:COUNT + %d`, inc))
			return nil
		}).
		OneTx(tx)
	if err != nil {
		return err
	}

	return nil
}
