package jdb

import (
	"fmt"

	"github.com/cgalvisleon/et/et"
)

var series *Model

/**
* defineSeries
* @param db *DB
* @return error
**/
func defineSeries(db *DB) error {
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
		"indexes":      []string{RECORDID},
	})
	if err != nil {
		return err
	}

	if err = series.Init(); err != nil {
		return err
	}

	return nil
}

/**
* GetSeries
* @param kind, tag string
* @return et.Item, error
**/
func GetSeries(kind, tag string) (et.Item, error) {
	if series == nil {
		return et.Item{}, fmt.Errorf(MSG_SERIES_NOT_DEFINED)
	}

	result, err := series.
		Where(Eq("kind", kind)).
		And(Eq("tag", tag)).
		One()
	if err != nil {
		return et.Item{}, err
	}

	return result, nil
}

/**
* GenSeries
* @param kind, tag string
* @return string, error
**/
func GenSeries(kind, tag string) (string, error) {
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
