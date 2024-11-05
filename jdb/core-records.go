package jdb

var Records *Model

func DefineRecords(db *Database) error {
	if Records != nil {
		return nil
	}

	if err := DefineSchemaCore(db); err != nil {
		return err
	}

	Records = NewModel(SchemaCore, "records")
	Records.DefineColumn(CreatedAtField, TypeDataTime, DefaultNow)
	Records.DefineColumn(UpdatedAtField, TypeDataTime, DefaultNow)
	Records.DefineColumn("table_schema", TypeDataText, DefaultNone)
	Records.DefineColumn("table_name", TypeDataText, DefaultNone)
	Records.DefineColumn("option", TypeDataShort, DefaultNone)
	Records.DefineColumn(SystemKeyField, TypeDataKey, DefaultKey)
	Records.DefineColumn(IndexField, TypeDataSerie, 0)
	Records.DefineKey("table_schema", "table_name", SystemKeyField)
	Records.DefineIndex(true, "table_schema", "table_name", "option", IndexField)
	Records.DefineFunction("RECORDS_BEFORE_INSERT", TpSqlFunction, `
		CREATE OR REPLACE FUNCTION core.RECORDS_BEFORE_INSERT()
		RETURNS
			TRIGGER AS $$  
		BEGIN
			IF NEW._IDT = '-1' THEN
				NEW._IDT = uuid_generate_v4();
			END IF;

			PERFORM pg_notify(
			'before',
			json_build_object(
				'schema', TG_TABLE_SCHEMA,
				'table', TG_TABLE_NAME,
				'option', TG_OP,        
				'_idt', NEW._IDT
			)::text
			);
		RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`)

	return nil
}
