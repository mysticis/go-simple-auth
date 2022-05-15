package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDatabase(t *testing.T) {

	err := InitDatabase()

	assert.NoError(t, err)
}
