package jdb

import (
	"encoding/json"
	"fmt"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/timezone"
)

var models *Model

/**
* defineModel
* @param db *Database
* @return error
**/
func defineModel(db *Database) error {
	if models != nil {
		return nil
	}

	var err error
	models, err = db.Define(et.Json{
		"schema":  "core",
		"name":    "models",
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
				"name": "name",
				"type": "text",
			},
			{
				"name": "version",
				"type": "int",
			},
			{
				"name": "definition",
				"type": "bytes",
			},
			{
				"name": RECORDID,
				"type": "key",
			},
		},
		"record_field": RECORDID,
		"primary_keys": []string{"kind", "name"},
		"indices":      []string{"version", RECORDID},
	})
	if err != nil {
		return err
	}

	if err = models.Init(); err != nil {
		return err
	}

	return nil
}

/**
* setModel
* @param kind string, name string, version int, definition []byte
* @return error
**/
func setModel(kind, name string, version int, definition []byte) error {
	if models == nil {
		return nil
	}

	now := timezone.Now()
	data := et.Json{
		"kind":       kind,
		"name":       name,
		"version":    version,
		"definition": definition,
	}
	_, err := models.
		Upsert(data).
		BeforeInsertOrUpdate(func(tx *Tx, data et.Json) error {
			data.Set("created_at", now)
			data.Set("updated_at", now)
			return nil
		}).
		BeforeUpdate(func(tx *Tx, data et.Json) error {
			data.Set("updated_at", now)
			return nil
		}).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

/**
* loadModel
* @param kind string, name string
* @return et.Json
**/
func loadModel(kind, name string, v any) error {
	if models == nil {
		return fmt.Errorf(MSG_MODEL_NOT_FOUND, name)
	}

	items, err := models.
		Query(et.Json{
			"where": et.Json{
				"kind": et.Json{
					"eq": kind,
				},
				"and": []et.Json{
					{
						"name": et.Json{
							"eq": name,
						},
					},
				},
			},
		}).
		One()
	if err != nil {
		return err
	}

	if !items.Ok {
		return nil
	}

	scr, err := items.Byte("definition")
	if err != nil {
		return err
	}

	err = json.Unmarshal(scr, v)
	if err != nil {
		return err
	}

	return nil
}

/**
* deleteModel
* @param name string
* @return error
**/
func deleteModel(name string) error {
	if models == nil {
		return nil
	}

	_, err := models.
		Delete(et.Json{
			"where": et.Json{
				"name": et.Json{
					"eq": name,
				},
			},
		}).
		Exec()
	if err != nil {
		return err
	}

	return nil
}
