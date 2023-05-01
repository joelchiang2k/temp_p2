package controllers

import "github.com/gin-gonic/gin"

func Authenticate(c *gin.Context) {

	var requestBody map[string]interface{}

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
	
	if username == "ece30861defaultadminuser" && isAdmin == true && password == "correcthorsebatterystaple123(!__+@**(A'\"`;DROP TABLE packages;" {
		c.String(200, "token")
	} else {
		c.JSON(401, "The user or password is invalid.")
	}
}
