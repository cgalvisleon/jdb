package oci

func (s *Oracle) CreateCore() error {
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
	if err := s.defineModel(); err != nil {
		return err
	}

	return nil
}
