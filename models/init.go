package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	//dsn := "postgres://postgres@localhost:5432/ben_temp"
	dsn := "host=34.171.210.8 user=postgres password=461dbpassword dbname=postgres port=5432 sslmode=disable"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = database.AutoMigrate(&PackageCreate{})
	if err != nil {
		return
	}

	err = database.AutoMigrate(&Token{})
	if err != nil {
		return
	}

	defaultUser := Token{
		Username:  "ece30861defaultadminuser",
		Password:  "correcthorsebatterystaple123(!__+@**(A'\"`;DROP TABLE packages;",
		AuthToken: "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	}

	var user Token

	check := database.Where("Username = ?", defaultUser.Username).First(&user)
	if check.Error == gorm.ErrRecordNotFound {
		database.Create(&defaultUser)
	}

	DB = database
}