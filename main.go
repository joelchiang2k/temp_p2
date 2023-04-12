package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//router.Use(static.Serve("/", static.LocalFile("./src", true)))

	api := router.Group("/package")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "hello",
			})
		})
	}

	api.GET("", CreatePackage)
	router.Run(":9876")
}

func CreatePackage(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{
		"message": "create package has not been implemented.",
	})
}
