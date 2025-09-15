package postgres

import (
	"fmt"

	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* GrantPrivileges
* @param username, database string
* @return error
**/
func (s *Postgres) GrantPrivileges(username, database string) error {
	/* Grant privileges */
	grantDatabase := fmt.Sprintf("GRANT CONNECT ON DATABASE %s TO %s;", database, username)
	err := jdb.Definition(s.jdb, grantDatabase)
	if err != nil {
		return err
	}

	/* Grant schema */
	grantSchema := fmt.Sprintf("GRANT USAGE ON SCHEMA public TO %s;", username)
	err = jdb.Definition(s.jdb, grantSchema)
	if err != nil {
		return err
	}

	/* Grant tables */
	grantTables := fmt.Sprintf("GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO %s;", username)
	err = jdb.Definition(s.jdb, grantTables)
	if err != nil {
		return err
	}

	/* Revoke drop */
	revokeDrop := fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE ALL ON TABLES FROM %s;", username)
	err = jdb.Definition(s.jdb, revokeDrop)
	if err != nil {
		return err
	}

	/* Grant future tables */
	grantFutureTables := fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO %s;", username)
	err = jdb.Definition(s.jdb, grantFutureTables)
	if err != nil {
		return err
	}

	return nil
}

/**
* CreateUser
* @param username, password, confirmation string
* @return error
**/
func (s *Postgres) CreateUser(username, password, confirmation string) error {
	if password != confirmation {
		return fmt.Errorf("password do not match!")
	}

	query := fmt.Sprintf("CREATE ROLE %s WITH LOGIN PASSWORD '%s';", username, password)
	err := jdb.Definition(s.jdb, query)
	if err != nil {
		return err
	}

	grantPrivilegesQuery := fmt.Sprintf(`GRANT ALL PRIVILEGES ON DATABASE %s;`, username)
	err = jdb.Definition(s.jdb, grantPrivilegesQuery)
	if err != nil {
		return err
	}

	return nil
}

/**
* ChangePassword
* @param username, password, confirmation string
* @return error
**/
func (s *Postgres) ChangePassword(username, password, confirmation string) error {
	if password != confirmation {
		return fmt.Errorf("password do not match!")
	}

	query := fmt.Sprintf("ALTER ROLE %s WITH PASSWORD '%s';", username, password)
	err := jdb.Definition(s.jdb, query)
	if err != nil {
		return err
	}

	return nil
}

/**
* DeleteUser
* @param username string
* @return error
**/
func (s *Postgres) DeleteUser(username string) error {
	query := fmt.Sprintf("DROP ROLE IF EXISTS %s;", username)
	err := jdb.Definition(s.jdb, query)
	if err != nil {
		return err
	}

	return nil
}
