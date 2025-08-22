package jdb

import (
	"fmt"
	"sync"

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

/**
* upsertTable
* @param tx *Tx, schema, name string, inc int
* @return error
**/
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

func (s *DB) Cardinality() error {
	if coreTables == nil || !coreTables.isInit {
		return nil
	}

	page := 1
	ok := true
	limit := 100
	for ok {
		items, err := coreTables.
			Select(
				"table_name",
			).
			Page(page).
			Rows(limit)
		if err != nil {
			return err
		}

		ok = items.Ok
		if !ok {
			break
		}

		wg := sync.WaitGroup{}
		for _, item := range items.Result {
			wg.Add(1)
			go func(item et.Json) {
				defer wg.Done()

				table := item.Str("table_name")
				model, err := LoadModel(s, table)
				if err != nil {
					console.Error(err)
					return
				}

				count, err := model.Counted()
				if err != nil {
					console.Error(err)
					return
				}

				_, err = coreTables.
					Update(et.Json{
						"count": count,
					}).
					Where("table_name").Eq(table).
					One()
				if err != nil {
					console.Error(err)
					return
				}
			}(item)
		}
		wg.Wait()

		page++
	}

	return nil
}
