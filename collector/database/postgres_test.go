package database_test

import (
	"collector/database"
	"collector/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPostgresDB(t *testing.T) {
	mockDB := &mocks.IDb{}
	db := database.NewPostgres(mockDB)

	_, ok := db.(*postgresDB)
	assert.True(t, ok)
}

func TestExecute(t *testing.T) {
	mockDB := new(mocks.IDb)
	mockDB.Mock.
}
