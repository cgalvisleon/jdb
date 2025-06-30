package mysql

import (
	"fmt"

	"github.com/cgalvisleon/et/mistake"
	jdb "github.com/cgalvisleon/jdb/jdb"
)

/**
* GrantPrivileges
* @param username, database string
* @return error
**/
func (s *Mysql) GrantPrivileges(username, database string) error {
	/* Grant privileges */
	grantDatabase := fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s;", database, username)
	_, err := jdb.Exec(s.db, grantDatabase)
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
func (s *Mysql) CreateUser(username, password, confirmation string) error {
	if password != confirmation {
		return mistake.New("password do not match!")
	}

	query := fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s';", username, password)
	_, err := jdb.Exec(s.db, query)
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
func (s *Mysql) ChangePassword(username, password, confirmation string) error {
	if password != confirmation {
		return mistake.New("password do not match!")
	}

	query := fmt.Sprintf("ALTER USER %s WITH PASSWORD '%s';", username, password)
	_, err := jdb.Exec(s.db, query)
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
func (s *Mysql) DeleteUser(username string) error {
	query := fmt.Sprintf("DROP USER IF EXISTS %s;", username)
	_, err := jdb.Exec(s.db, query)
	if err != nil {
		return err
	}

	return nil
}
