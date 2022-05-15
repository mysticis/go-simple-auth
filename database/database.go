package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// global DB object to be used accross packages

var GlobalDB *gorm.DB

func InitDatabase() (err error) {

	dbURL := "postgres://postgres:secret@localhost:5432/test2"

	GlobalDB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)

	}

	return
}
