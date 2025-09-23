package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
)

var series *Model

/**
* defineSeries
* @param db *Database
* @return error
**/
func defineSeries(db *Database) error {
	if series != nil {
		return nil
	}

	var err error
	series, err = db.Define(et.Json{
		"schema":  "core",
		"name":    "series",
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
				"name": "kind",
				"type": "key",
			},
			{
				"name": "tag",
				"type": "key",
			},
			{
				"name": "value",
				"type": "int",
			},
			{
				"name": "format",
				"type": "key",
			},
			{
				"name": RECORDID,
				"type": "key",
			},
		},
		"record_field": RECORDID,
		"primary_keys": []string{"kind", "tag"},
		"indices":      []string{RECORDID},
		"debug":        true,
	})
	if err != nil {
		return err
	}

	series.IsCore = true
	err = series.Init()
	if err != nil {
		return err
	}

	return nil
}

/**
* GetSeries
* @param kind, tag string
* @return string, error
**/
func GetSeries(kind, tag string) (string, error) {
	if series == nil {
		return "", fmt.Errorf(MSG_SERIES_NOT_DEFINED)
	}

	item, err := series.
		Upsert(et.Json{
			"kind": kind,
			"tag":  tag,
		}).
		BeforeInsert(func(tx *Tx, data et.Json) error {
			data.Set("format", "%08d")
			data.Set("value", 1)
			return nil
		}).
		BeforeUpdate(func(tx *Tx, data et.Json) error {
			data.Set("value", "value + 1")
			return nil
		}).
		Where(Eq("kind", kind)).
		And(Eq("tag", tag)).
		Return(et.Json{
			"value":  "value",
			"format": "format",
		}).
		One()
	if err != nil {
		return "", err
	}

	value := item.Int("value")
	format := item.Str("format")
	result := fmt.Sprintf(format, value)

	return result, nil
}

/**
* SetSeries
* @param kind, tag, format string
* @return error
**/
func SetSeries(kind, tag, format string, lastValue int) error {
	if series == nil {
		return fmt.Errorf(MSG_SERIES_NOT_DEFINED)
	}

	_, err := series.
		Upsert(et.Json{
			"kind":   kind,
			"tag":    tag,
			"format": format,
			"value":  lastValue,
		}).
		Exec()
	if err != nil {
		return err
	}

	return nil
}
