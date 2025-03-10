package sqlite

import (
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/jdb/jdb"
)

func (s *SqlLite) SetFlow(name string, value []byte) error
func (s *SqlLite) GetFlow(id string) (jdb.Flow, error)
func (s *SqlLite) DeleteFlow(id string) error
func (s *SqlLite) FindFlows(search string, page, rows int) (et.List, error)
