package controllers

import (
	"ex/part2/models"

	"github.com/gin-gonic/gin"
)

/*type psqlRow struct {
	Name string `json:"Name"`
	Content string `json:"Content"`
	ID string `json:"id"`
	Version string `json:"Version"`
	URL string `json:"url"`
}*/

type psqlRow struct {
	Name string `json:"Name"`
	Version string `json:"Version"`
}

type RegexString struct {
	RegEx string `json:"RegEx"`
}

func ByRegex(c *gin.Context) {
	var newRegEx RegexString

	if err := c.BindJSON(&newRegEx); err != nil {
		c.JSON(400, "bad request")
		return
	}

	var matchedPackages []psqlRow

	//var pVersion string
	//rows, err := models.DB.Raw("SELECT * from package_creates WHERE name = ?;", newRegEx.RegEx).Rows()
	rows, err := models.DB.Raw("SELECT * from package_creates WHERE name ~* ?;", newRegEx.RegEx).Rows()
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
	//	models.DB.ScanRows(rows, &matchedPackages)
		var match models.PackageCreate 
		//rows.Scan(&pName, &pVersion)
		err := rows.Scan(&match.Name, &match.Content, &match.ID, &match.Version, &match.URL)
		if err != nil {
			panic(err)
		}
		if err = rows.Err(); err != nil {
			panic(err)
		}
		var nameVersionFromRow psqlRow
		nameVersionFromRow.Name = match.Name
		nameVersionFromRow.Version = match.Version
		//match.Version = pVersion
		matchedPackages = append(matchedPackages, nameVersionFromRow)
	}

<<<<<<< HEAD
	
	//if rows empty no package found w/ regex string -> 404
	if(rows == nil){
		c.JSON(404, "No package found under this regex.")
	}

	c.JSON(200, gin.H {
		"value": matchedPackages,
	})


=======
>>>>>>> 08e4fa7 (updating metrics)
}
