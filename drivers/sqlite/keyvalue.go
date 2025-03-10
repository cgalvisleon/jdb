package sqlite

import "github.com/cgalvisleon/et/et"

func (s *SqlLite) SetKey(key string, value []byte) error
func (s *SqlLite) GetKey(key string) (et.KeyValue, error)
func (s *SqlLite) DeleteKey(key string) error
func (s *SqlLite) FindKeys(search string, page, rows int) (et.List, error)
