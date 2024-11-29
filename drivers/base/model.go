package base

import jdb "github.com/cgalvisleon/jdb/jdb"

func (s *Base) CreateModel(model *jdb.Model) error {
	return nil
}

func (s *Base) MutateModel(model *jdb.Model) error {
	return nil
}

func (s *Base) DefaultValue(tp jdb.TypeData) interface{}
