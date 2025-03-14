package postgres

import (
	"github.com/cgalvisleon/et/strs"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) CreateCore() error {
	if err := s.defineRecords(); err != nil {
		return err
	}
	if err := s.defineSeries(); err != nil {
		return err
	}
	if err := s.defineRecycling(); err != nil {
		return err
	}
	if err := s.defineDDL(); err != nil {
		return err
	}
	if err := s.defineModel(); err != nil {
		return err
	}
	if err := s.defineFunctions(s.nodeId); err != nil {
		return err
	}
	if err := s.defineKeyValue(); err != nil {
		return err
	}
	if err := s.defineFlows(); err != nil {
		return err
	}
	if err := s.defineCache(); err != nil {
		return err
	}

	return nil
}

func (s *Postgres) existTable(schema, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.tables
		WHERE UPPER(table_schema) = UPPER($1)
		AND UPPER(table_name) = UPPER($2));`

	items, err := s.Query(sql, schema, name)
	if err != nil {
		return false, err
	}

	if items.Count == 0 {
		return false, nil
	}

	return items.Bool(0, "exists"), nil
}

func parceSQL(sql string) string {
	return strs.Change(sql,
		[]string{"date_make", "date_update", "_id", "_idt", "_state"},
		[]string{jdb.CREATED_AT, jdb.UPDATED_AT, jdb.PRIMARYKEY, jdb.SYSID, jdb.STATUS})
}
