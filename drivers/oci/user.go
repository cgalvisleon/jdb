package oci

import (
	"fmt"

	"github.com/cgalvisleon/et/mistake"
)

func (s *Postgres) GrantPrivileges(username, database string) error {
	// Otorgar acceso a la base de datos
	grantDatabase := fmt.Sprintf("GRANT CONNECT ON DATABASE %s TO %s;", database, username)
	_, err := s.db.Exec(grantDatabase)
	if err != nil {
		return err
	}

	// Otorgar acceso al esquema público
	grantSchema := fmt.Sprintf("GRANT USAGE ON SCHEMA public TO %s;", username)
	_, err = s.db.Exec(grantSchema)
	if err != nil {
		return err
	}

	// Otorgar permisos de lectura y escritura en todas las tablas actuales
	grantTables := fmt.Sprintf("GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO %s;", username)
	_, err = s.db.Exec(grantTables)
	if err != nil {
		return err
	}

	// Asegurar que el usuario no pueda eliminar tablas ni modificar el esquema
	revokeDrop := fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE ALL ON TABLES FROM %s;", username)
	_, err = s.db.Exec(revokeDrop)
	if err != nil {
		return err
	}

	// Establecer que los permisos sean aplicables a futuras tablas creadas
	grantFutureTables := fmt.Sprintf("ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO %s;", username)
	_, err = s.db.Exec(grantFutureTables)
	if err != nil {
		return err
	}

	return nil
}

func (s *Postgres) CreateUser(username, password, confirmation string) error {
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
