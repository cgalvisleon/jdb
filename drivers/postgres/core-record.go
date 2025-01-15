package postgres

import (
	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) defineRecords() error {
	exist, err := s.existTable("core", "RECORDS")
	if err != nil {
		return console.Panic(err)
	}

	if exist {
		return s.defineRecordsFunction()
	}

	sql := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE EXTENSION IF NOT EXISTS pgcrypto;
	CREATE SCHEMA IF NOT EXISTS core;

  CREATE TABLE IF NOT EXISTS core.RECORDS(
    DATE_MAKE TIMESTAMP DEFAULT NOW(),
		DATE_UPDATE TIMESTAMP DEFAULT NOW(),
    TABLE_SCHEMA VARCHAR(80) DEFAULT '',
    TABLE_NAME VARCHAR(80) DEFAULT '',
		OPTION VARCHAR(80) DEFAULT '',
		SYNC BOOLEAN DEFAULT FALSE,
    _IDT VARCHAR(80) DEFAULT '-1',
    INDEX SERIAL,
    PRIMARY KEY (TABLE_SCHEMA, TABLE_NAME, _IDT)
  );    
  CREATE INDEX IF NOT EXISTS RECORDS_TABLE_SCHEMA_IDX ON core.RECORDS(TABLE_SCHEMA);
  CREATE INDEX IF NOT EXISTS RECORDS_TABLE_NAME_IDX ON core.RECORDS(TABLE_NAME);
	CREATE INDEX IF NOT EXISTS RECORDS_OPTION_IDX ON core.RECORDS(OPTION);
	CREATE INDEX IF NOT EXISTS RECORDS_SYNC_IDX ON core.RECORDS(SYNC);
  CREATE INDEX IF NOT EXISTS RECORDS__IDT_IDX ON core.RECORDS(_IDT);  
	CREATE INDEX IF NOT EXISTS RECORDS_INDEX_IDX ON core.RECORDS(INDEX);`

	err = s.Exec(sql)
	if err != nil {
		return console.Panic(err)
	}

	return s.defineRecordsFunction()
}

func (s *Postgres) defineRecordsFunction() error {
	sql := `
	CREATE OR REPLACE FUNCTION core.SYNC_NOTIFY()
  RETURNS
    TRIGGER AS $$  
  BEGIN
    IF NEW.SYNC = FALSE THEN      
			PERFORM pg_notify(
			'sync',
			json_build_object(
				'schema', NEW.TABLE_SCHEMA,
				'table', NEW.TABLE_NAME,
				'option', NEW.OPTION,        
				'_idt', NEW._IDT
			)::text
			);
		END IF;
  RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

	DROP TRIGGER IF EXISTS SYNC_AFTER_INSERT ON core.RECORDS CASCADE;
	CREATE TRIGGER SYNC_AFTER_INSERT
	AFTER INSERT ON core.RECORDS
	FOR EACH ROW
	EXECUTE PROCEDURE core.SYNC_NOTIFY();

	DROP TRIGGER IF EXISTS SYNC_AFTER_UPDATE ON core.RECORDS CASCADE;
	CREATE TRIGGER SYNC_AFTER_UPDATE
	AFTER UPDATE ON core.RECORDS
	FOR EACH ROW
	EXECUTE PROCEDURE core.SYNC_NOTIFY();

	CREATE OR REPLACE FUNCTION core.RECORDS_BEFORE_INSERT()
  RETURNS
    TRIGGER AS $$
	DECLARE
		VSYNC BOOLEAN;
  BEGIN
    IF NEW._IDT = '-1' THEN
      NEW._IDT = uuid_generate_v4();
			VSYNC = FALSE;
		ELSE
			VSYNC = TRUE;
		END IF;

		INSERT INTO core.RECORDS(TABLE_SCHEMA, TABLE_NAME, OPTION, SYNC, _IDT)
		VALUES (TG_TABLE_SCHEMA, TG_TABLE_NAME, TG_OP, VSYNC, NEW._IDT);

		PERFORM pg_notify(
		'trigger',
		json_build_object(
			'schema', TG_TABLE_SCHEMA,
			'table', TG_TABLE_NAME,
			'option', TG_OP,
			'event', 'before_insert',
			'_idt', NEW._IDT
		)::text
		);

  	RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION core.RECORDS_AFTER_INSERT()
  RETURNS
    TRIGGER AS $$
  BEGIN
		PERFORM pg_notify(
		'trigger',
		json_build_object(
			'schema', TG_TABLE_SCHEMA,
			'table', TG_TABLE_NAME,
			'option', TG_OP,
			'event', 'after_insert',
			'_idt', NEW._IDT
		)::text
		);
    
  	RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION core.RECORDS_BEFORE_UPDATE()
  RETURNS
    TRIGGER AS $$
	DECLARE
		VSYNC BOOLEAN;
  BEGIN
		IF OLD._IDT = NEW._IDT THEN
			VSYNC = FALSE;
		ELSE
			NEW._IDT = OLD._IDT;
			VSYNC = TRUE;
		END IF;

		UPDATE core.RECORDS SET
		DATE_UPDATE=NOW(),
		OPTION=TG_OP,
		SYNC=VSYNC
		WHERE _IDT=NEW._IDT;

		PERFORM pg_notify(
		'trigger',
		json_build_object(
			'schema', TG_TABLE_SCHEMA,
			'table', TG_TABLE_NAME,
			'option', TG_OP,
			'event', 'before_update',
			'_idt', NEW._IDT
		)::text
		);

  	RETURN NEW;
  END;
  $$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION core.RECORDS_AFTER_UPDATE()
	RETURNS
		TRIGGER AS $$
	BEGIN
		PERFORM pg_notify(
		'trigger',
		json_build_object(
			'schema', TG_TABLE_SCHEMA,
			'table', TG_TABLE_NAME,
			'option', TG_OP,
			'event', 'after_update',
			'_idt', NEW._IDT
		)::text
		);

		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

  CREATE OR REPLACE FUNCTION core.RECORDS_BEFORE_DELETE()
  RETURNS
    TRIGGER AS $$  
  BEGIN
		UPDATE core.RECORDS SET
		DATE_UPDATE=NOW(),
		OPTION=TG_OP,
		SYNC=FALSE
		WHERE _IDT=OLD._IDT;
		
		PERFORM pg_notify(
		'trigger',
		json_build_object(
			'schema', TG_TABLE_SCHEMA,
			'table', TG_TABLE_NAME,
			'option', TG_OP,
			'event', 'before_delete',
			'_idt', OLD._IDT
		)::text
		);

  	RETURN OLD;
  END;
  $$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION core.RECORDS_AFTER_DELETE()
	RETURNS
		TRIGGER AS $$
	BEGIN
	PERFORM pg_notify(
		'trigger',
		json_build_object(
			'schema', TG_TABLE_SCHEMA,
			'table', TG_TABLE_NAME,
			'option', TG_OP,
			'event', 'after_delete',
			'_idt', NEW._IDT
		)::text
		);

	RETURN OLD;
	END;
	$$ LANGUAGE plpgsql;
	`
	err := s.Exec(sql)
	if err != nil {
		return console.Panic(err)
	}

	return nil
}

func defineRecordTrigger(table string) string {
	result := jdb.SQLDDL(`
	DROP TRIGGER IF EXISTS RECORDS_BEFORE_INSERT ON $1 CASCADE;
	CREATE TRIGGER RECORDS_BEFORE_INSERT
	BEFORE INSERT ON $1
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_BEFORE_INSERT();

	DROP TRIGGER IF EXISTS RECORDS_AFTER_INSERT ON $1 CASCADE;
	CREATE TRIGGER RECORDS_AFTER_INSERT
	BEFORE INSERT ON $1
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_AFTER_INSERT();

	DROP TRIGGER IF EXISTS RECORDS_BEFORE_UPDATE ON $1 CASCADE;
	CREATE TRIGGER RECORDS_BEFORE_UPDATE
	BEFORE UPDATE ON $1
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_BEFORE_UPDATE();

	DROP TRIGGER IF EXISTS RECORDS_AFTER_UPDATE ON $1 CASCADE;
	CREATE TRIGGER RECORDS_AFTER_UPDATE
	BEFORE UPDATE ON $1
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_AFTER_UPDATE();

	DROP TRIGGER IF EXISTS RECORDS_BEFORE_DELETE ON $1 CASCADE;
	CREATE TRIGGER RECORDS_BEFORE_DELETE
	BEFORE DELETE ON $1
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_BEFORE_DELETE();

	DROP TRIGGER IF EXISTS RECORDS_AFTER_DELETE ON $1 CASCADE;
	CREATE TRIGGER RECORDS_AFTER_DELETE
	BEFORE DELETE ON $1
	FOR EACH ROW
	EXECUTE PROCEDURE core.RECORDS_AFTER_DELETE();`, table)

	result = strs.Replace(result, "\t", "")

	return result
}
