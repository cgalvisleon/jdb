package sqlite

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
)

func (s *SqlLite) Exec(sql string, arg ...any) error
func (s *SqlLite) Query(sql string, arg ...any) (et.Items, error)
func (s *SqlLite) One(sql string, arg ...any) (et.Item, error)
func (s *SqlLite) Data(source, sql string, arg ...any) (et.Items, error)
func (s *SqlLite) Select(ql *jdb.Ql) (et.Items, error)
func (s *SqlLite) Count(ql *jdb.Ql) (int, error)
func (s *SqlLite) Exists(ql *jdb.Ql) (bool, error)
