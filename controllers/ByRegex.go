package controllers

import (
	"ex/part2/models"

	"github.com/gin-gonic/gin"
)

type RegexString struct {
	RegEx string `json:"RegEx"`
}
func ByRegex(c *gin.Context) {
	var newRegEx RegexString

	if err := c.BindJSON(&newRegEx); err != nil {
		c.JSON(400, "bad request")
	}

	rows, err := models.DB.Raw("SELECT * from package_creates WHERE name ~ 'nodist'").Rows()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	//if rows empty no package found w/ regex string -> 404

	
}