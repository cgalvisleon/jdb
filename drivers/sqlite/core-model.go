package sqlite

import "github.com/cgalvisleon/et/console"

func (s *SqlLite) defineModel() error {
	sql := parceSQL(`	
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
	CREATE INDEX IF NOT EXISTS RECORDS_DATE_MAKE_IDX ON core.RECORDS(DATE_MAKE);
	CREATE INDEX IF NOT EXISTS RECORDS_DATE_UPDATE_IDX ON core.RECORDS(DATE_UPDATE);
  CREATE INDEX IF NOT EXISTS RECORDS_TABLE_SCHEMA_IDX ON core.RECORDS(TABLE_SCHEMA);
  CREATE INDEX IF NOT EXISTS RECORDS_TABLE_NAME_IDX ON core.RECORDS(TABLE_NAME);
	CREATE INDEX IF NOT EXISTS RECORDS_OPTION_IDX ON core.RECORDS(OPTION);
	CREATE INDEX IF NOT EXISTS RECORDS_SYNC_IDX ON core.RECORDS(SYNC);
  CREATE INDEX IF NOT EXISTS RECORDS__IDT_IDX ON core.RECORDS(_IDT);  
	CREATE INDEX IF NOT EXISTS RECORDS_INDEX_IDX ON core.RECORDS(INDEX);`)

	err := s.Exec(sql)
	if err != nil {
		return console.Panic(err)
	}

	return nil
}
