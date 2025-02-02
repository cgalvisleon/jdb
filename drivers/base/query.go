package base

import (
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Base) Exec(sql string, params ...any) error {
	return nil
}

func (s *Base) Query(sql string, params ...any) (et.Items, error) {
	return et.Items{}, nil
}

func (s *Base) Data(source, sql string, params ...any) (et.Items, error) {
	return et.Items{}, nil
}

func (s *Base) Select(ql *jdb.Ql) (et.Items, error) {
	return et.Items{}, nil
}

func (s *Base) Count(ql *jdb.Ql) (int, error) {
	return 0, nil
}
