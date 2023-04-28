package controllers

import (
	"github.com/gin-gonic/gin"
)
func HandleAuth(c *gin.Context){
	c.JSON(501, "This system does not support authentication.")
}