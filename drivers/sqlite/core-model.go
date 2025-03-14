package sqlite

import "github.com/cgalvisleon/et/console"

func (s *SqlLite) defineModel() error {
	sql := parceSQL(`	
  `)

	err := s.Exec(sql)
	if err != nil {
		return console.Panic(err)
	}

	return nil
}
