package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	dbURL := "postgres://postgres@localhost:5432/ben_temp"

	database, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = database.AutoMigrate(&PackageCreate{})
	if err != nil {
		return
	}

	DB = database
}