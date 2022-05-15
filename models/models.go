package models

import (
	"github.com/mysticis/go-net/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty" gorm:"unique"`
	Password string `json:"password,omitempty"`
}

//createUserRecord creates a user in the database

func (user *User) CreateUserRecord() error {

	result := database.GlobalDB.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

//hashPassword encrypts the user password

func (user *User) HashPassword(password string) error {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return err
	}

	user.Password = string(bytes)

	return nil

}

func (user *User) CheckPassword(providedPassword string) error {

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))

	if err != nil {
		return err
	}

	return nil
}
