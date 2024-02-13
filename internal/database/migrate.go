package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattes/migrate/source/file"
)

func RunMigrate(user, pass, host, port, database string) error {
	getwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, database)
	configPath := "file://" + filepath.Join(getwd, "migrations")

	m, err := migrate.New(
		configPath,
		dsn,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	err = m.Up()
	if err != nil {
		log.Println(err)
		return err
	}
	m.Close()

	return nil
}
