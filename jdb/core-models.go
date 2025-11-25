package jdb

import (
	"encoding/json"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/timezone"
)

var models *Model

/**
* defineModel
* @param db *DB
* @return error
**/
func defineModel(db *DB) error {
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
		"primary_keys": []string{"name"},
		"indexes":      []string{"version", RECORDID},
	})
	if err != nil {
		return err
	}

	models.isCore = true
	if err = models.Init(); err != nil {
		return err
	}

	return nil
}

/**
* setModel
* @param name string, version int, definition []byte
* @return error
**/
func setModel(name string, version int, definition []byte) error {
	if models == nil {
		return nil
	}

	now := timezone.Now()
	data := et.Json{
		"name":       name,
		"version":    version,
		"definition": definition,
	}
	_, err := models.
		Upsert(data).
		BeforeInsertOrUpdate(func(tx *Tx, old, new et.Json) error {
			new.Set("created_at", now)
			new.Set("updated_at", now)
			return nil
		}).
		BeforeUpdate(func(tx *Tx, old, new et.Json) error {
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
* @param name string
* @return et.Json
**/
func loadModel(name string, v any) error {
	if models == nil {
		return ErrModelNotFound
	}

	items, err := models.
		Where(Eq("A.name", name)).
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
		Delete().
		Where(Eq("name", name)).
		Exec()
	if err != nil {
		return err
	}

	return nil
}
