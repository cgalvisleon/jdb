package sqlite

import (
	"fmt"

	"github.com/cgalvisleon/jdb/jdb"
)

func table(schema, name string) string {
	return fmt.Sprintf(`%s_%s`, schema, name)
}

func tableByModel(model *jdb.Model) string {
	return table(model.Schema, model.Name)
}

func (s *SqlLite) LoadModel(model *jdb.Model) error {
	return nil
}

func (s *SqlLite) DropModel(model *jdb.Model) error {
	return nil
}

func (s *SqlLite) MutateModel(model *jdb.Model) error {
	return nil
}
