package controllers

import (
	"ex/part2/models"
	"fmt"

	"gorm.io/gorm"
)

func HandleAuth(authHeader string) int {
	var token models.Token
	var output int
	// var count int64
	// xAuthHeader := c.Request.Header["X-Authorization"]
	if len(authHeader) == 0 {
		fmt.Println("No X-Authorization header found")
		output = 400
		return output
	}

	// models.DB.Model(&token{}).Where("AuthToken = ?", authHeader).Count(&count)

	// if err := models.DB.Where("AuthToken = ?", authHeader).First(&token).Error; err != gorm.ErrRecordNotFound {
	// 	fmt.Println("Token found")
	// 	output = 0
	// } else {
	// 	fmt.Println("Token not found")
	// 	output = 400
	// }

	err := models.DB.Where("AuthToken = ?", authHeader).First(&token).Error
	if err != gorm.ErrRecordNotFound {
    	fmt.Println("Token found")
		fmt.Println(token)
    	output = 0
	} else {
		fmt.Println("Token not found")
		output = 400
	}
	return output
}
