package postgres

import (
	jdb "github.com/cgalvisl/jdb/pkg"
	"github.com/cgalvisleon/et/logs"
	"github.com/cgalvisleon/et/strs"
)

func (s *Postgres) defineSeries() error {
	exist, err := s.existTable("core", "SERIES")
	if err != nil {
		return logs.Panic(err)
	}

	if exist {
		return s.defineSeriesFunction()
	}

	sql := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE EXTENSION IF NOT EXISTS pgcrypto;
	CREATE SCHEMA IF NOT EXISTS core;

  CREATE TABLE IF NOT EXISTS core.SERIES(
		SERIE VARCHAR(250) DEFAULT '',
		VALUE BIGINT DEFAULT 0,
		PRIMARY KEY(SERIE)
	);`

	err = s.Exec(sql)
	if err != nil {
		return logs.Panic(err)
	}

	return s.defineSeriesFunction()
}

func (s *Postgres) defineSeriesFunction() error {
	sql := `
	CREATE OR REPLACE FUNCTION core.nextserie(tag VARCHAR(250))
	RETURNS BIGINT AS $$
	DECLARE
	 result BIGINT;
	BEGIN
	 UPDATE core.SERIES SET
	 VALUE = VALUE + 1
	 WHERE SERIE = tag
	 RETURNING VALUE INTO result;
	 IF NOT FOUND THEN
	  INSERT INTO core.SERIES(SERIE, VALUE)
		VALUES (tag, 1)
		RETURNING VALUE INTO result;
	 END IF;

	 RETURN COALESCE(result, 0);
	END;
	$$ LANGUAGE plpgsql;
	
	CREATE OR REPLACE FUNCTION core.setserie(tag VARCHAR(250), val BIGINT)
	RETURNS BIGINT AS $$
	DECLARE
	 result BIGINT;
	BEGIN
	 UPDATE core.SERIES SET
	 VALUE = val
	 WHERE SERIE = tag
	 RETURNING VALUE INTO result;
	 IF NOT FOUND THEN
	  INSERT INTO core.SERIES(SERIE, VALUE)
		VALUES (tag, val)
		RETURNING VALUE INTO result;	
	 END IF;

	 RETURN COALESCE(result, 0);
	END;
	$$ LANGUAGE plpgsql;
	
	CREATE OR REPLACE FUNCTION core.currserie(tag VARCHAR(250))
	RETURNS BIGINT AS $$
	DECLARE
	 result BIGINT;
	BEGIN
	 SELECT VALUE INTO result
	 FROM core.SERIES
	 WHERE SERIE = tag LIMIT 1;

	 RETURN COALESCE(result, 0);
	END;
	$$ LANGUAGE plpgsql;
	
	CREATE OR REPLACE FUNCTION core.SERIES_AFTER_SET()
  RETURNS
    TRIGGER AS $$
	DECLARE
		TAG VARCHAR(250);
  BEGIN
	  SELECT CONCAT(TG_TABLE_SCHEMA, '.',  TG_TABLE_NAME) INTO TAG;
		PERFORM core.setserie(TAG, NEW.INDEX);

  	RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;
	`

	err := s.Exec(sql)
	if err != nil {
		return logs.Panic(err)
	}

	return nil
}

func (s *Postgres) GetSerie(tag string) int64 {
	db := s.db
	if s.master != nil {
		db = s.master
	}

	sql := `SELECT core.nextserie($1) AS SERIE;`
	rows, err := s.query(db, sql, tag)
	if err != nil {
		logs.Error(err)
		return 0
	}
	defer rows.Close()

	item := jdb.RowsToItem(rows)
	if !item.Ok {
		return 0
	}

	result := item.Int64("serie")

	return result
}

func (s *Postgres) SetSerie(tag string, val int) int64 {
	db := s.db
	if s.master != nil {
		db = s.master
	}

	sql := `SELECT core.setserie($1) AS SERIE;`
	rows, err := s.query(db, sql, tag)
	if err != nil {
		logs.Error(err)
		return 0
	}
	defer rows.Close()

	item := jdb.RowsToItem(rows)
	if !item.Ok {
		return 0
	}

	result := item.Int64("serie")

	return result
}

func (s *Postgres) CurrentSerie(tag string) int64 {
	db := s.db
	if s.master != nil {
		db = s.master
	}

	sql := `SELECT core.currserie($1) AS SERIE;`
	rows, err := s.query(db, sql, tag)
	if err != nil {
		logs.Error(err)
		return 0
	}
	defer rows.Close()

	item := jdb.RowsToItem(rows)
	if !item.Ok {
		return 0
	}

	result := item.Int64("serie")

	return result
}

func (s *Postgres) NextCode(tag, prefix string) string {
	num := s.GetSerie(tag)

	if len(prefix) == 0 {
		return strs.Format("%08v", num)
	} else {
		return strs.Format("%s%08v", prefix, num)
	}
}
