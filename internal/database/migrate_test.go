package database_test

import (
	"api/internal/database"
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

func Test_RunMigrate(t *testing.T) {
	pool, err := dockertest.NewPool("")
	assert.Nil(t, err)
	assert.NotNil(t, pool)

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	assert.Nil(t, err)

	host := "localhost"
	user := "db"
	pass := "db"
	name := "postgres"

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("postgres", "latest", []string{"POSTGRES_USER=" + user, "POSTGRES_PASSWORD=" + pass})
	assert.Nil(t, err)
	assert.NotNil(t, resource)

	port := resource.GetPort("5432/tcp")

	var db *sql.DB

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)
		db, err = sql.Open("pgx", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	err = database.RunMigrate(user, pass, host, port, name)
	assert.Nil(t, err)

	err = pool.Purge(resource)
	assert.Nil(t, err)
}
