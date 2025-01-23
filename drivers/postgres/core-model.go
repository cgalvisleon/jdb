package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* defineModel create alter or delete
* @return error
**/
func (s *Postgres) defineModel() error {
	exist, err := s.existTable("core", "MODELS")
	if err != nil {
		return console.Panic(err)
	}

	if exist {
		return nil
	}

	sql := parceSQL(`
  CREATE TABLE IF NOT EXISTS core.MODELS(
		TABLENAME VARCHAR(80) DEFAULT '',
		VERSION INTEGER DEFAULT 0,
		MODEL BYTEA,
		INDEX SERIAL,
		PRIMARY KEY(TABLENAME)
	);
	CREATE INDEX IF NOT EXISTS MODELS_TABLENAME_IDX ON core.MODELS(TABLENAME);
	CREATE INDEX IF NOT EXISTS MODELS_INDEX_IDX ON core.MODELS(INDEX);`)

	err = s.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

/**
* getModel
* @params table string
* @return et.Item, error
**/
func (s *Postgres) getModel(table string) (et.Item, error) {
	sql := `
	SELECT *
	FROM core.MODELS
	WHERE TABLENAME = $1
	LIMIT 1;`

	return s.One(jdb.Select, sql, table)
}

/**
* upsertModel
* @params table string
* @params version int
* @params model []byte
* @return error
**/
func (s *Postgres) upsertModel(table string, version int, model []byte) error {
	sql := `
	UPDATE core.MODELS SET
	MODEL = $3,
	VERSION = $2		
	WHERE TABLENAME = $1
	RETURNING *;`

	item, err := s.One(jdb.Select, sql, table, version, model)
	if err != nil {
		return err
	}

	if item.Ok {
		return nil
	}

	sql = `
	INSERT INTO core.MODELS(TABLENAME, MODEL, VERSION)
	VALUES ($1, $2, $3);`

	err = s.Exec(sql, table, model, version)
	if err != nil {
		console.Alertm(et.Json{
			"table":   table,
			"version": version,
			"error":   err.Error(),
		}.ToString())
		return err
	}

	return nil
}

/**
* deleteModel
* @params table string
* @return error
**/
func (s *Postgres) deleteModel(table string) error {
	sql := `
	DELETE FROM core.MODELS
	WHERE TABLENAME = $1;`

	err := s.Exec(sql, table)
	if err != nil {
		return err
	}

	return nil
}
