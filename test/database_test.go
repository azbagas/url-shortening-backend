package test

import (
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	db := SetupTestDB()
	defer db.Close()

	err := db.Ping()
	assert.Nil(t, err)
}
