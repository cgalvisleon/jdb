package base

import (
	jdb "github.com/cgalvisl/jdb/pkg"
	"github.com/cgalvisleon/et/et"
)

func (s *Base) Exec(sql string, params ...interface{}) error {
	return nil
}

func (s *Base) SQL(sql string, params ...interface{}) (et.Items, error) {
	return et.Items{}, nil
}

func (s *Base) Query(linq *jdb.Linq) (et.Items, error) {
	return et.Items{}, nil
}

func (s *Base) One(sql string, params ...interface{}) (et.Item, error) {
	return et.Item{}, nil
}

func (s *Base) Count(linq *jdb.Linq) (int, error) {
	return 0, nil
}

func (s *Base) Last(linq *jdb.Linq) (et.Items, error) {
	return et.Items{}, nil
}
