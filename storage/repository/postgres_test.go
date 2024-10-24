// postgres_test.go
package repository_test

import (
	"telegram_reminder_bot/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgresDB(t *testing.T) {
	cfg := repository.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "user",
		Password: "password",
		DBName:   "testdb",
		SSLMode:  "disable",
		TIMEZONE: "UTC",
	}

	db, err := repository.NewPostgresDB(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Close the database after the test is done to prevent resource leaks.
	defer db.Close()
}
