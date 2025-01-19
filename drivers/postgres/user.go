package postgres

import (
	"fmt"

	"github.com/cgalvisleon/et/mistake"
)

func (s *Postgres) CreateUser(username, password, confirmation string) error {
	if password != confirmation {
		return mistake.New("password do not match!")
	}

	query := fmt.Sprintf("CREATE ROLE %s WITH LOGIN PASSWORD '%s';", username, password)
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (s *Postgres) ChangePassword(username, password, confirmation string) error {
	if password != confirmation {
		return mistake.New("password do not match!")
	}

	query := fmt.Sprintf("ALTER ROLE %s WITH PASSWORD '%s';", username, password)
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (s *Postgres) DeleteUser(username string) error {
	query := fmt.Sprintf("DROP ROLE IF EXISTS %s;", username)
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
