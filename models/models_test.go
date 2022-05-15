package models

import (
	"os"
	"testing"

	"github.com/mysticis/go-net/database"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	user := User{
		Password: "secretPassword",
	}

	err := user.HashPassword(user.Password)

	assert.NoError(t, err)
	os.Setenv("passwordHash", user.Password)
}

func TestCreateUserRecord(t *testing.T) {

	var userResult User

	err := database.InitDatabase()

	if err != nil {
		t.Error(err)
	}

	err = database.GlobalDB.AutoMigrate(&User{})

	assert.NoError(t, err)

	user := User{
		Name:     "test user",
		Email:    "test@email.com",
		Password: os.Getenv("passwordHash"),
	}

	err = user.CreateUserRecord()

	assert.NoError(t, err)

	database.GlobalDB.Where("email = ?", user.Email).Find(&userResult)

	database.GlobalDB.Unscoped().Delete(&user)

	assert.Equal(t, "test user", userResult.Name)

	assert.Equal(t, "test@email.com", userResult.Email)

}

func TestCheckPassword(t *testing.T) {

	hash := os.Getenv("passwordHash")

	user := User{
		Password: hash,
	}

	err := user.CheckPassword("secretPassword")

	assert.NoError(t, err)
}
