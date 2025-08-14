package jdb

/**
* GetType
* @param name string
* @return interface{}
**/
func (s *Model) GetType(name string) interface{} {
	result, ok := s.Types[name]
	if !ok {
		return []string{}
	}

	return result
}

/**
* SetType
* @param name string, tp interface{}
**/
func (s *Model) SetType(name string, tp interface{}) error {
	s.Types[name] = tp
	if s.isInit {
		return s.Save()
	}

	return nil
}
