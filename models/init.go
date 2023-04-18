package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	dsn := "postgres://postgres@localhost:5432/ben_temp"
	//dsn := "host=35.224.74.124 user=postgres password=461dbpassword dbname=postgres port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = database.AutoMigrate(&PackageCreate{})
	if err != nil {
		return
	}

	DB = database
}