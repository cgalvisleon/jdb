package oracle

import (
	"fmt"

	"github.com/cgalvisleon/et/mistake"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* grantPrivileges
* @param username, database string
* @return error
**/
func (s *Oracle) grantPrivileges(username, database string) error {
	/* Grant privileges */
	grantDatabase := fmt.Sprintf("GRANT CONNECT ON DATABASE %s TO %s;", database, username)
	_, err := jdb.Exec(s.db, grantDatabase)
	if err != nil {
		return err
	}

	/* Grant schema */
	grantSchema := fmt.Sprintf("GRANT USAGE ON SCHEMA public TO %s;", username)
	_, err = jdb.Exec(s.db, grantSchema)
	if err != nil {
		return err
	}

	/* Grant tables */
	grantTables := fmt.Sprintf("GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO %s;", username)
	_, err = jdb.Exec(s.db, grantTables)
	if err != nil {
		return err
	}

	/* Revoke drop */
	revokeDrop := fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE ALL ON TABLES FROM %s;", username)
	_, err = jdb.Exec(s.db, revokeDrop)
	if err != nil {
		return err
	}

	/* Grant future tables */
	grantFutureTables := fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO %s;", username)
	_, err = jdb.Exec(s.db, grantFutureTables)
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
func (s *Oracle) createUser(username, password, confirmation string) error {
	if password != confirmation {
		return mistake.New("password do not match!")
	}

	query := fmt.Sprintf("CREATE ROLE %s WITH LOGIN PASSWORD '%s';", username, password)
	_, err := jdb.Exec(s.db, query)
	if err != nil {
		return err
	}

	grantPrivilegesQuery := fmt.Sprintf(`GRANT ALL PRIVILEGES ON DATABASE %s;`, username)
	_, err = jdb.Exec(s.db, grantPrivilegesQuery)
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
func (s *Oracle) changePassword(username, password, confirmation string) error {
	if password != confirmation {
		return mistake.New("password do not match!")
	}

	query := fmt.Sprintf("ALTER ROLE %s WITH PASSWORD '%s';", username, password)
	_, err := jdb.Exec(s.db, query)
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
func (s *Oracle) deleteUser(username string) error {
	query := fmt.Sprintf("DROP ROLE IF EXISTS %s;", username)
	_, err := jdb.Exec(s.db, query)
	if err != nil {
		return err
	}

	return nil
}
