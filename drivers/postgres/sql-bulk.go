package postgres

import (
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Postgres) sqlBulk(command *jdb.Command) string {
	return s.sqlInsert(command)
}
