package base

import "github.com/cgalvisleon/et/et"

func (s *Base) CreateDatabase(name string) error {
	return nil
}

func (s *Base) DropDatabase(name string) error {
	return nil
}

func (s *Base) RenameDatabase(name, newname string) error {
	return nil
}

func (s *Base) SetParams(data et.Json) error {
	return nil
}
