package jdb

func (s *Command) inserted() error {
	err := s.bulk()
	if err != nil {
		return err
	}

	return nil
}
