package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) defineRecycling() error {
	exist, err := s.existTable("core", "RECYCLING")
	if err != nil {
		return console.Panic(err)
	}

	if exist {
		return s.defineRecyclingFunction()
	}

	sql := strs.Change(`
  CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE EXTENSION IF NOT EXISTS pgcrypto;
	CREATE SCHEMA IF NOT EXISTS core;
  
  CREATE TABLE IF NOT EXISTS core.RECYCLING(
    DATE_MAKE TIMESTAMP DEFAULT NOW(),
    TABLE_SCHEMA VARCHAR(80) DEFAULT '',
    TABLE_NAME VARCHAR(80) DEFAULT '',
    _IDT VARCHAR(80) DEFAULT '-1',
    INDEX SERIAL,
    PRIMARY KEY (TABLE_SCHEMA, TABLE_NAME, _IDT)
  );    
  CREATE INDEX IF NOT EXISTS RECYCLING_TABLE_SCHEMA_IDX ON core.RECYCLING(TABLE_SCHEMA);
  CREATE INDEX IF NOT EXISTS RECYCLING_TABLE_NAME_IDX ON core.RECYCLING(TABLE_NAME);
  CREATE INDEX IF NOT EXISTS RECYCLING__IDT_IDX ON core.RECYCLING(_IDT);
	CREATE INDEX IF NOT EXISTS RECYCLING_INDEX_IDX ON core.RECYCLING(INDEX);`,
		[]string{"date_create", "date_update", "_id", "_idt", "_data"},
		[]string{jdb.CreatedAtField.Str(), jdb.UpdatedAtField.Str(), jdb.KeyField.Str(), jdb.SystemKeyField.Str(), jdb.SourceField.Str()})

	err = s.Exec(sql)
	if err != nil {
		return console.Panic(err)
	}

	return s.defineRecyclingFunction()
}

func (s *Postgres) defineRecyclingFunction() error {
	sql := strs.Change(`  
  CREATE OR REPLACE FUNCTION core.RECYCLING_UPDATE()
  RETURNS
    TRIGGER AS $$
  DECLARE
    CHANNEL VARCHAR(250);
  BEGIN
    IF NEW._STATE != OLD._STATE AND NEW._STATE = '-2' THEN      
      INSERT INTO core.RECYCLING(TABLE_SCHEMA, TABLE_NAME, _IDT)
      VALUES (TG_TABLE_SCHEMA, TG_TABLE_NAME, NEW._IDT);

      PERFORM pg_notify(
      'recycling',
      json_build_object(
        'schema', TG_TABLE_SCHEMA,
        'table', TG_TABLE_NAME,
        '_idt', NEW._IDT
      )::text
      );
		ELSEIF NEW._STATE != OLD._STATE THEN
      DELETE FROM core.RECYCLING
      WHERE TABLE_SCHEMA = TG_TABLE_SCHEMA
      AND TABLE_NAME = TG_TABLE_NAME
      AND _IDT = NEW._IDT;
    END IF;

  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;
  
  CREATE OR REPLACE FUNCTION core.RECYCLING_DELETE()
  RETURNS
    TRIGGER AS $$  
  BEGIN
    DELETE FROM core.RECYCLING
    WHERE TABLE_SCHEMA = TG_TABLE_SCHEMA
    AND TABLE_NAME = TG_TABLE_NAME
    AND _IDT = OLD._IDT;

  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;`,
		[]string{"date_create", "date_update", "_id", "_idt", "_data"},
		[]string{jdb.CreatedAtField.Str(), jdb.UpdatedAtField.Str(), jdb.KeyField.Str(), jdb.SystemKeyField.Str(), jdb.SourceField.Str()})

	err := s.Exec(sql)
	if err != nil {
		return console.Panic(err)
	}

	return nil
}

func defineRecyclingTrigger(table string) string {
	result := strs.Change(`
  DROP TRIGGER IF EXISTS RECYCLING ON $1 CASCADE;
	CREATE TRIGGER RECYCLING
	AFTER UPDATE ON $1
	FOR EACH ROW WHEN (OLD._STATE!=NEW._STATE)
	EXECUTE PROCEDURE core.RECYCLING_UPDATE();

	DROP TRIGGER IF EXISTS RECYCLING_DELETE ON $1 CASCADE;
	CREATE TRIGGER RECYCLING_DELETE
	AFTER DELETE ON $1
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECYCLING_DELETE();`,
		[]string{"_STATE", "$1"},
		[]string{jdb.StateField.Str(), table})

	result = strs.Replace(result, "\t", "")

	return result
}
