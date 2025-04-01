package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
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
    TABLENAME TEXT DEFAULT '',
    VERSION INTEGER DEFAULT 0,
    MODEL BYTEA,
    _IDT VARCHAR(80) DEFAULT '-1',
		INDEX BIGINT DEFAULT 0,
		PRIMARY KEY(TABLENAME)
	);
	CREATE INDEX IF NOT EXISTS MODELS_INDEX_IDX ON core.MODELS(INDEX);`)
	sql = strs.Append(sql, defineRecordTrigger("core.DDL"), "\n")
	sql = strs.Append(sql, defineSeriesTrigger("core.DDL"), "\n")

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
