package postgres

func (s *Postgres) CreateCore() error {
	if err := s.defineRecords(); err != nil {
		return err
	}
	if err := s.defineSeries(); err != nil {
		return err
	}
	if err := s.defineRecycling(); err != nil {
		return err
	}
	if err := s.defineDDL(); err != nil {
		return err
	}

	return nil
}

func (s *Postgres) existTable(schema, name string) (bool, error) {
	sql := `
	SELECT EXISTS(
		SELECT 1
		FROM information_schema.tables
		WHERE UPPER(table_schema) = UPPER($1)
		AND UPPER(table_name) = UPPER($2));`

	item, err := s.One(sql, schema, name)
	if err != nil {
		return false, err
	}

	return item.Bool("exists"), nil
}
