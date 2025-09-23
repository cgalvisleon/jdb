package jdb

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/timezone"
)

var records *Model

/**
* defineRecord
* @param db *Database
* @return error
**/
func defineRecord(db *Database) error {
	if records != nil {
		return nil
	}

	var err error
	records, err = db.Define(et.Json{
		"schema":  "core",
		"name":    "records",
		"version": 1,
		"columns": []et.Json{
			{
				"name": "created_at",
				"type": "datetime",
			},
			{
				"name": "updated_at",
				"type": "datetime",
			},
			{
				"name": "schema",
				"type": "key",
			},
			{
				"name": "table",
				"type": "text",
			},
			{
				"name": RECORDID,
				"type": "key",
			},
		},
		"record_field": RECORDID,
		"primary_keys": []string{"schema", "table"},
		"indices":      []string{RECORDID},
		"debug":        true,
	})
	if err != nil {
		return err
	}

	records.isCore = true
	records.BeforeInsert(func(tx *Tx, data et.Json) error {
		data.Set("created_at", timezone.Now())
		data.Set("updated_at", timezone.Now())
		return nil
	})
	records.BeforeUpdate(func(tx *Tx, data et.Json) error {
		data.Set("updated_at", timezone.Now())
		return nil
	})
	records.BeforeDelete(func(tx *Tx, data et.Json) error {
		data.Set("updated_at", timezone.Now())
		return nil
	})

	err = records.Init()
	if err != nil {
		return err
	}

	return nil
}

/**
* upsertRecord
* @param schema, table string
* @return error
**/
func upsertRecord(schema, table string) error {
	if records == nil {
		return nil
	}

	_, err := records.
		Upsert(et.Json{
			"schema": schema,
			"table":  table,
		}).
		AfterInsert(func(tx *Tx, data et.Json) error {
			data.Set("created_at", timezone.Now())
			data.Set("updated_at", timezone.Now())
			return nil
		}).
		AfterUpdate(func(tx *Tx, data et.Json) error {
			data.Set("updated_at", timezone.Now())
			return nil
		}).
		AfterDelete(func(tx *Tx, data et.Json) error {
			data.Set("updated_at", timezone.Now())
			return nil
		}).
		Exec()
	if err != nil {
		return err
	}

	return nil
}
