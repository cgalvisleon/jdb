package sqlite

import "github.com/cgalvisleon/jdb/jdb"

func (s *SqlLite) existTable(schema, name string) (bool, error) {
	table := table(schema, name)
	sql := `
	SELECT name
	FROM sqlite_master
	WHERE type='table'
	AND name=?;`

	items, err := jdb.Query(s.db, sql, table)
	if err != nil {
		return false, err
	}

	if items.Count == 0 {
		return false, nil
	}

	return items.Bool(0, "exists"), nil
}
