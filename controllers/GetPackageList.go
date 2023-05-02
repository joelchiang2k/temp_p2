package controllers

import (
	"encoding/json"
	"ex/part2/models"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type packagesMetadataStruct struct {
	ID string `json:"ID,omitempty"`
	Name string `json:"Name"`
	Version string `json:"Version"`
}

type PackageQueryArray struct {
	QueryArr []packagesMetadataStruct `json:""`
}

//var packagesToReturn []PackagesMetadataStruct
func GetPackageList(c *gin.Context) {
	var newQuery []packagesMetadataStruct

	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err!= nil {
		c.JSON(400, "SThere is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
		return
	}


	err = json.Unmarshal(reqBody, &newQuery)
	if err!= nil {
		c.JSON(400, "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
		return
	}
	fmt.Println(newQuery)

	var packagesToReturn []packagesMetadataStruct
	//var foundPackage PackagesMetadataStruct
	for _, PackagesMetadataStruct := range newQuery{

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