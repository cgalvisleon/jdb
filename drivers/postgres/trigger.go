package postgres

import "database/sql"

func triggerBeforeInsert(db *sql.DB) error {
	sql := `
	CREATE SCHEMA IF NOT EXISTS core;
	CREATE TABLE IF NOT EXISTS core.tables (
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	schema VARCHAR(80),
	table TEXT,
	total BIGINT);
	ALTER TABLE core.tables ADD CONSTRAINT tables_pk PRIMARY KEY (schema, table);

	CREATE OR REPLACE FUNCTION core.before_insert_records()
	RETURNS TRIGGER AS $$	
	BEGIN
		IF TD_OP = 'INSERT' THEN
			INSERT INTO core.tables (created_at, updated_at, schema, table, total)
			VALUES (now(), now(), TG_TABLE_SCHEMA, TG_TABLE_NAME, 1)
			ON CONFLICT (schema, table) DO UPDATE
			SET total = core.tables.total + 1,
					updated_at = now();
		END IF;

		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION core.before_delete_records()
	RETURNS TRIGGER AS $$	
	BEGIN
		IF TD_OP = 'DELETE' THEN
			UPDATE core.tables
			SET total = core.tables.total - 1,
					updated_at = now()
			WHERE schema = TG_TABLE_SCHEMA
			AND table = TG_TABLE_NAME;
		END IF;

		RETURN OLD;
	END;
	$$ LANGUAGE plpgsql;
	`
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
