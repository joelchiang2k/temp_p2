package controllers

import (
	"encoding/json"
	"ex/part2/models"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type MetadataStruct struct {
	Name string `json:"Name"`
	Version string `json:"Version"`
	ID uint `json:"ID"`
}

type PackageDataStruct struct {
	Content string `json:"Content"`
	URL string `json:"URL"`
}

type FullPackageData struct {
	Metadata MetadataStruct `json:"metadata"`
	Data PackageDataStruct `json:"data"`
}

func UpdatePackage(c *gin.Context) {
	var pkg models.PackageCreate
	var packageToUpdate FullPackageData

	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	reqBodyJson := string(reqBody)
	fmt.Println("request before trying to bind for UPDATE")
	fmt.Println(reqBodyJson)
	//validate request
	if err := c.BindJSON(&packageToUpdate); err != nil {
		c.JSON(400, "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
		return
	}
	//fmt.Println(packageToUpdate)
	niceJSON, err := json.MarshalIndent(packageToUpdate, "", " ")
	if err != nil {
		fmt.Println(err)
		return
	}	
	fmt.Println("AFTER BINDING JSON")
	fmt.Println("package/byid PUT request")
	fmt.Println(string(niceJSON))

	//get package if exists in db
	//incorrect format in route string
	if c.Param("{id}") == "/" {
		c.JSON(400, "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
		return
	} else if err := models.DB.Where("id = ?", c.Param("{id}")).First(&pkg).Error; err != nil {
		//package not found -> return error 404
		c.JSON(404, "Package does not exist.")
		return
	}


	//validate name, ID, version are matched
	if(packageToUpdate.Metadata.Name == pkg.Name && packageToUpdate.Metadata.ID == pkg.ID && packageToUpdate.Metadata.Version == pkg.Version){
		if(packageToUpdate.Data.Content != ""){
			//models.DB.Update("content", packageToUpdate.Data.Content)
			models.DB.Table("package_creates").Where("id = ?", c.Param("{id}")).Updates(map[string]interface{}{"content": packageToUpdate.Data.Content})
			c.JSON(200, "Package was updated.")
			return
		}else if(packageToUpdate.Data.URL != ""){
			//operate on url which fucking blows

		}
	}else{
		c.JSON(404, "ID, Name, and Version do not all match.")
		return
	}

}