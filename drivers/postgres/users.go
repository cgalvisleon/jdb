package postgres

import (
	"fmt"

	"github.com/cgalvisleon/et/mistake"
)

/**
* grantPrivileges
* @param username, database string
* @return error
**/
func (s *Postgres) grantPrivileges(username, database string) error {
	/* Grant privileges */
	grantDatabase := fmt.Sprintf("GRANT CONNECT ON DATABASE %s TO %s;", database, username)
	_, err := s.db.Exec(grantDatabase)
	if err != nil {
		return err
	}

	/* Grant schema */
	grantSchema := fmt.Sprintf("GRANT USAGE ON SCHEMA public TO %s;", username)
	_, err = s.db.Exec(grantSchema)
	if err != nil {
		return err
	}

	/* Grant tables */
	grantTables := fmt.Sprintf("GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO %s;", username)
	_, err = s.db.Exec(grantTables)
	if err != nil {
		return err
	}

	/* Revoke drop */
	revokeDrop := fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE ALL ON TABLES FROM %s;", username)
	_, err = s.db.Exec(revokeDrop)
	if err != nil {
		return err
	}

	/* Grant future tables */
	grantFutureTables := fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO %s;", username)
	_, err = s.db.Exec(grantFutureTables)
	if err != nil {
		return err
	}

	return nil
}

/**
* createUser
* @param username, password, confirmation string
* @return error
**/
func (s *Postgres) createUser(username, password, confirmation string) error {
	if password != confirmation {
		return mistake.New("password do not match!")
	}

	query := fmt.Sprintf("CREATE ROLE %s WITH LOGIN PASSWORD '%s';", username, password)
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	grantPrivilegesQuery := fmt.Sprintf(`GRANT ALL PRIVILEGES ON DATABASE %s;`, username)
	_, err = s.db.Exec(grantPrivilegesQuery)
	if err != nil {
		return err
	}

	return nil
}

/**
* changePassword
* @param username, password, confirmation string
* @return error
**/
func (s *Postgres) changePassword(username, password, confirmation string) error {
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

/**
* deleteUser
* @param username string
* @return error
**/
func (s *Postgres) deleteUser(username string) error {
	query := fmt.Sprintf("DROP ROLE IF EXISTS %s;", username)
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
