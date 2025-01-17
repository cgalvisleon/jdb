package postgres

/**
* defineFuntion: Define functions if not exists
* @return error
**/
func (s *Postgres) defineFunctions() error {
	sql := `
	CREATE OR REPLACE FUNCTION core.add_constraint_if_not_exists(
    schema_name TEXT,
    table_name TEXT,
    constraint_name TEXT,
    constraint_definition TEXT
	) RETURNS VOID AS $$
	BEGIN
		IF NOT EXISTS (
				SELECT 1
				FROM pg_catalog.pg_constraint c
				JOIN pg_catalog.pg_namespace n ON n.oid = c.connamespace
				JOIN pg_catalog.pg_class t ON t.oid = c.conrelid
				WHERE n.nspname = schema_name
					AND t.relname = table_name
					AND c.conname = constraint_name
		) THEN
				EXECUTE constraint_definition;
		END IF;
	END;
	$$ LANGUAGE plpgsql;`

	err := s.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
