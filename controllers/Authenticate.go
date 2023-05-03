package controllers

import (
	"ex/part2/models"
	// "fmt"

	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {

	var requestBody map[string]interface{}
	var foundToken models.Token

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, ok := requestBody["User"].(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{"error": "There is missing field(s) in the AuthenticationRequest or it is formed improperly."})
		return
	}

	username, ok := user["name"]
	if !ok {
		c.JSON(400, gin.H{"error": "There is missing field(s) in the AuthenticationRequest or it is formed improperly."})
		return
	}

	isAdmin, ok := user["isAdmin"]
	if !ok {
		c.JSON(400, gin.H{"error": "There is missing field(s) in the AuthenticationRequest or it is formed improperly."})
		return
	}

	secret, ok := requestBody["Secret"].(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{"error": "There is missing field(s) in the AuthenticationRequest or it is formed improperly."})
		return
	}

	password, ok := secret["password"]
	if !ok {
		c.JSON(400, gin.H{"error": "There is missing field(s) in the AuthenticationRequest or it is formed improperly."})
		return
	}

	if isAdmin == true {
		if err := models.DB.Where("username = ? AND password = ?", username, password).First(&foundToken).Error; err != nil {
			c.JSON(404, "The user or password is invalid.")
		} else {
			c.JSON(200, foundToken.Auth_token)
		}
	} else {
		c.JSON(401, "The user or password is invalid.")
	}
}
