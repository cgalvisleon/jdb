package jdb

import "github.com/cgalvisleon/et/et"

var recycling *Model

/**
* defineRecycling
* @param db *Database
* @return error
**/
func defineRecycling(db *Database) error {
	if recycling != nil {
		return nil
	}

	var err error
	recycling, err = db.Define(et.Json{
		"schema":  "core",
		"name":    "recycling",
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
		"primary_keys": []string{RECORDID},
		"indices":      []string{},
		"debug":        true,
	})
	if err != nil {
		return err
	}

	recycling.IsCore = true
	err = recycling.Init()
	if err != nil {
		return err
	}

	return nil
}

/**
* SetRecycling
* @param schema, table, index string
* @return error
**/
func (s *Database) SetRecycling(schema, table, index string) error {
	_, err := recycling.
		Upsert(et.Json{
			"schema": schema,
			"table":  table,
			RECORDID: index,
		}).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

/**
* DeleteRecycling
* @param schema, table, index string
* @return error
**/
func (s *Database) DeleteRecycling(schema, table, index string) error {
	_, err := recycling.
		Delete(et.Json{
			"where": et.Json{
				"schema": schema,
				"table":  table,
				RECORDID: index,
			},
		}).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

/**
* QueryRecycling
* @param query et.Json
* @return (et.Items, error)
**/
func (s *Database) QueryRecycling(query et.Json) (et.Items, error) {
	items, err := recycling.
		Query(query).
		All()
	if err != nil {
		return et.Items{}, err
	}

	return items, nil
}
