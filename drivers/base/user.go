package base

func (s *Base) GrantPrivileges(username, dbName string) error {
	return nil
}

func (s *Base) CreateUser(username, password, confirmation string) error {
	return nil
}

func (s *Base) ChangePassword(username, password, confirmation string) error {
	return nil
}

func (s *Base) DeleteUser(username string) error {
	return nil
}
