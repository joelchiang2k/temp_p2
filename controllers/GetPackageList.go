package controllers

import (
	"encoding/json"
	"ex/part2/models"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type packagesMetadataStruct struct {
	ID      string `json:"ID,omitempty"`
	Name    string `json:"Name"`
	Version string `json:"Version"`
}

type PackageQueryArray struct {
	QueryArr []packagesMetadataStruct `json:""`
}

//var packagesToReturn []PackagesMetadataStruct
func GetPackageList(c *gin.Context) {
	var newQuery []packagesMetadataStruct
	var token models.Token

	authHeader := c.Request.Header.Get("X-Authorization")
	if len(authHeader) == 0 {
		c.JSON(400, "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
		return
	}
	if err := models.DB.Where("Auth_token = ?", authHeader).First(&token).Error; err != nil {
		fmt.Println("Token not found")
		c.JSON(400, "There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid.")
	}

	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
		return
	}

	err = json.Unmarshal(reqBody, &newQuery)
	if err != nil {
		c.JSON(400, "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
		return
	}
	fmt.Println(newQuery)

	var packagesToReturn []packagesMetadataStruct
	//var foundPackage PackagesMetadataStruct
	for _, PackagesMetadataStruct := range newQuery {

		var foundPackage packagesMetadataStruct
		if err := models.DB.Table("package_creates").Where("name = ?", PackagesMetadataStruct.Name).Find(&foundPackage).Error; err != nil {
			c.JSON(400, "package not found")
			return
		}
		packagesToReturn = append(packagesToReturn, foundPackage)
	}
	fmt.Println(packagesToReturn)
	c.JSON(200, packagesToReturn)
	return
}
