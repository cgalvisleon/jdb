package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/utility"
)

/**
* defineDDL create alter or delete
* @return error
**/
func (s *Postgres) defineDDL() error {
	exist, err := s.existTable("core", "DDL")
	if err != nil {
		return console.Panic(err)
	}

	if exist {
		return nil
	}

	sql := parceSQL(`
  CREATE TABLE IF NOT EXISTS core.DDL(
		_ID VARCHAR(80) DEFAULT '-1',
		SQL BYTEA,
		_IDT VARCHAR(80) DEFAULT '-1',
		INDEX BIGINT DEFAULT 0,
		PRIMARY KEY(_ID)
	);
	CREATE INDEX IF NOT EXISTS DDL__IDT_IDX ON core.DDL(_IDT);
	CREATE INDEX IF NOT EXISTS DDL_INDEX_IDX ON core.DDL(INDEX);`)
	sql = strs.Append(sql, defineRecordTrigger("core.DDL"), "\n")
	sql = strs.Append(sql, defineSeriesTrigger("core.DDL"), "\n")

	err = s.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

/**
* upsertDDL
* @params query string
**/
func (s *Postgres) upsertDDL(id string, query string) error {
	sql := parceSQL(`
	UPDATE core.DDL SET
	SQL = $2
	WHERE _ID = $1
	RETURNING _ID;`)

	item, err := s.One(sql, id, []byte(query))
	if err != nil {
		return err
	}

	if item.Ok {
		return nil
	}

	sql = parceSQL(`
	INSERT INTO core.DDL(_ID, SQL, INDEX)
	VALUES ($1, $2, $3);`)

	id = utility.GenKey(id)
	index := s.GetSerie("ddl")
	err = s.Exec(sql, id, []byte(query), index)
	if err != nil {
		console.Alertm(et.Json{
			"_id": id,
			"sql": query,
		}.ToString())
		return err
	}

	return nil
}

/**
* deleteDDL
* @params query string
**/
func (s *Postgres) deleteDDL(id string) error {
	sql := parceSQL(`
	DELETE FROM core.DDL
	WHERE _ID = $1;`)

	err := s.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}
