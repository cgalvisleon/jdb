package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
)

var records *Model

/**
* defineRecords
* @param db *DB
* @return error
**/
func defineRecords(db *DB) error {
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
				"name": "table_schema",
				"type": "key",
			},
			{
				"name": "table_name",
				"type": "key",
			},
			{
				"name": RECORDID,
				"type": "key",
			},
		},
		"record_field": RECORDID,
		"primary_keys": []string{"table_schema", "table_name", RECORDID},
		"indexes":      []string{"updated_at"},
	})
	if err != nil {
		return err
	}

	records.isCore = true
	if err = records.Init(); err != nil {
		return err
	}

	return nil
}

/**
* GetRecord
* @param id string
* @return et.Item, error
**/
func GetRecordById(id string) (et.Item, error) {
	if records == nil {
		return et.Item{}, fmt.Errorf(MSG_RECORDS_NOT_DEFINED)
	}

	result, err := records.
		Where(Eq(RECORDID, id)).
		One()
	if err != nil {
		return et.Item{}, err
	}

	if !result.Ok {
		return et.Item{}, fmt.Errorf(MSG_RECORD_NOT_FOUND, id)
	}

	name := result.Str("name")
	model, err := GetModel(records.Database, name)
	if err != nil {
		return et.Item{}, err
	}

	return model.GetRecordById(id)
}
