package jdb

func (s *Command) deleted() error {
	if err := s.prepare(); err != nil {
		return err
	}

	model := s.From
	results, err := s.Db.Command(s)
	if err != nil {
		return err
	}

	s.Result = results
	s.ResultMap, err = model.getMapResultByPk(s.Result.Result)
	if err != nil {
		return err
	}

	for _, before := range s.ResultMap {
		for _, fn := range model.afterDelete {
			err := fn(s.tx, before)
			if err != nil {
				return err
			}
		}

		for _, fn := range s.afterDelete {
			err := fn(s.tx, before)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
