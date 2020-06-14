package migrations

import (
	"database/sql"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/squadcast_assignment/internal/config"
)

func Up(dbConfig config.Database) error {
	dbInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
	db, err := sql.Open(dbConfig.Dialect, dbInfo)

	if err != nil {
		return fmt.Errorf("connection to MySQL failed: %s", err.Error())
	}

	migrations := &migrate.FileMigrationSource{
		Dir: dbConfig.MigrationsDir,
	}

	_, err = migrate.Exec(db, dbConfig.Dialect, migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("migrations up failed due to: %v", err)
	}
	return nil
}

func Down(dbConfig config.Database) error {
	dbInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
	db, err := sql.Open(dbConfig.Dialect, dbInfo)

	if err != nil {
		return fmt.Errorf("connection to MySQL failed: %s", err.Error())
	}

	migrations := &migrate.FileMigrationSource{
		Dir: dbConfig.MigrationsDir,
	}

	_, err = migrate.Exec(db, dbConfig.Dialect, migrations, migrate.Down)
	if err != nil {
		return fmt.Errorf("migrations down failed due to: %v", err)
	}
	return nil
}
