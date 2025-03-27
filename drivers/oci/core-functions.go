package oci

import jdb "github.com/cgalvisleon/jdb/jdb"

/**
* defineFuntion: Define functions if not exists
* @return error
**/
func (s *Postgres) defineFunctions(nodeId int) error {
	sql := jdb.SQLDDL(`
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
	$$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION generate_snowflake_id()
	RETURNS BIGINT AS $$
	DECLARE
			epoch BIGINT := 1672531200000; -- Epoch base (2023-01-01 00:00:00 UTC)
			timestamp BIGINT;
			last_timestamp BIGINT := 0; -- Último timestamp registrado
			node_id INT := $1; -- Identificador único para el nodo (cambiar según el nodo)
			sequence INT := 0; -- Secuencia de IDs en el mismo milisegundo
			max_sequence INT := 4095; -- Límite de la secuencia (12 bits)
	BEGIN
			-- Obtén el timestamp actual en milisegundos
			timestamp := (EXTRACT(EPOCH FROM clock_timestamp()) * 1000)::BIGINT - epoch;

			-- Si el reloj retrocede, lanza un error
			IF timestamp < last_timestamp THEN
					RAISE EXCEPTION 'El reloj del sistema retrocedió.';
			END IF;

			-- Si el mismo timestamp, incrementa la secuencia
			IF timestamp = last_timestamp THEN
					sequence := sequence + 1;
					-- Si la secuencia alcanza su límite, espera al próximo milisegundo
					IF sequence > max_sequence THEN
							LOOP
									timestamp := (EXTRACT(EPOCH FROM clock_timestamp()) * 1000)::BIGINT - epoch;
									IF timestamp > last_timestamp THEN
											sequence := 0;
											EXIT;
									END IF;
							END LOOP;
					END IF;
			ELSE
					-- Reinicia la secuencia si el timestamp es nuevo
					sequence := 0;
			END IF;

			-- Actualiza el último timestamp registrado
			last_timestamp := timestamp;

			-- Genera el Snowflake ID
			RETURN (timestamp << 22) | (node_id << 12) | sequence;
	END;
	$$ LANGUAGE plpgsql;`, nodeId)

	err := s.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
