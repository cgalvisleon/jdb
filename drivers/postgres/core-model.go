package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
)

/**
* defineModel create alter or delete
* @return error
**/
func (s *Postgres) defineModel() error {
	sql := parceSQL(`
  CREATE TABLE IF NOT EXISTS core_MODELS (
    TABLENAME TEXT DEFAULT '',
    VERSION INTEGER DEFAULT 0,
    MODEL BLOB,
    INDEX_ID INTEGER PRIMARY KEY AUTOINCREMENT
	);

	CREATE INDEX IF NOT EXISTS MODELS_TABLENAME_IDX ON core_MODELS(TABLENAME);
	CREATE INDEX IF NOT EXISTS MODELS_INDEX_IDX ON core_MODELS(INDEX_ID);`)

	err := s.Exec(sql)
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

	result, err := s.Query(sql, table)
	if err != nil {
		return et.Item{}, err
	}

	return result.First(), nil
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

	item, err := s.Query(sql, table, version, model)
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

	sql = s.ddlTableDrop(table)
	err = s.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
