package base

import jdb "github.com/cgalvisleon/jdb/jdb"

func (s *Base) LoadTable(model *jdb.Model) (bool, error) {
	return false, nil
}

func (s *Base) LoadModel(model *jdb.Model) error {
	return nil
}

func (s *Base) DropModel(model *jdb.Model) error {
	return nil
}
