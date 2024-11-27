package josefina

import (
	jdb "github.com/cgalvisl/jdb/pkg"
	"github.com/cgalvisleon/et/et"
)

func (s *Josefina) Exec(sql string, params ...interface{}) error {
	return nil
}

func (s *Josefina) SQL(sql string, params ...interface{}) (et.Items, error) {
	return et.Items{}, nil
}

func (s *Josefina) Query(linq *jdb.Linq) (et.Items, error) {
	return et.Items{}, nil
}

func (s *Josefina) One(sql string, params ...interface{}) (et.Item, error) {
	return et.Item{}, nil
}

func (s *Josefina) Count(linq *jdb.Linq) (int, error) {
	return 0, nil
}

func (s *Josefina) Last(linq *jdb.Linq) (et.Items, error) {
	return et.Items{}, nil
}
