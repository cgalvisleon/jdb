package postgres

import jdb "github.com/cgalvisleon/jdb/jdb"

func (s *Postgres) CreateModel(model *jdb.Model) error {
	return nil
}

func (s *Postgres) MutateModel(model *jdb.Model) error {
	return nil
}

func (s *Postgres) DefaultValue(tp jdb.TypeData) interface{}
