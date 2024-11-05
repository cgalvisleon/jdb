package jdb

var SchemaCore *Schema

func DefineSchemaCore(db *Database) error {
	if SchemaCore != nil {
		return nil
	}

	SchemaCore = NewSchema(db, "core")
	if err := SchemaCore.Init(); err != nil {
		return err
	}

	return nil
}

func InitCore(db *Database) {
	if err := DefineRecords(db); err != nil {
		panic(err)
	}
}
