package base

import (
	"github.com/cgalvisleon/et/et"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

func (s *Base) Exec(sql string, params ...any) error {
	return nil
}

func (s *Base) SQL(sql string, params ...any) (et.Items, error) {
	return et.Items{}, nil
}

func (s *Base) Query(linq *jdb.Ql) (et.Items, error) {
	return et.Items{}, nil
}

func (s *Base) One(sql string, params ...any) (et.Item, error) {
	return et.Item{}, nil
}

func (s *Base) Count(linq *jdb.Ql) (int, error) {
	return 0, nil
}

func (s *Base) Last(linq *jdb.Ql) (et.Items, error) {
	return et.Items{}, nil
}
