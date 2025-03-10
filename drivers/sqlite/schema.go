package sqlite

import "github.com/cgalvisleon/jdb/jdb"

var schema *jdb.Schema
var _db *jdb.DB

func defineSchema(db *jdb.DB) error {
	if schema != nil {
		return nil
	}

	var err error
	_db = db
	schema, err = jdb.NewSchema(db, "core")
	if err != nil {
		return err
	}

	return nil
}

// Schema
func (s *SqlLite) CreateSchema(name string) error
func (s *SqlLite) DropSchema(name string) error
