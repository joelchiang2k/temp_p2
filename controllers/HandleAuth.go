package controllers

import (
	"ex/part2/models"
	"fmt"
)

func HandleAuth(authHeader string) int {
	var token models.Token
	var output int

	if len(authHeader) == 0 {
		fmt.Println("No X-Authorization header found")
		output = 400
		return output
	}

	if err := models.DB.Where("Auth_token = ?", authHeader).First(&token).Error; err != nil {
		output = 400
	} else {
		output = 0
	}

	return output
}
