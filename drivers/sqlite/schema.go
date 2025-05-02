package sqlite

import (
	"github.com/cgalvisleon/jdb/jdb"
)

var schema *jdb.Schema
var _db *jdb.DB

/* Schema */
func (s *SqlLite) LoadSchema(name string) error
func (s *SqlLite) DropSchema(name string) error
