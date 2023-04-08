package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type JSON_LIST struct {
	Version string `json:"Version"`
	Name string `json:"Name"`
}

func main() {
	router := gin.Default()

	//router.Use(static.Serve("/", static.LocalFile("./src", true)))

	api := router.Group("/package")
	{
		/*api.GET("/", func(c *gin.Context){
			c.JSON(http.StatusOK, gin.H{
				"message": "hello",
			})
		})*/
		api.GET("", CreatePackage)
		api.GET("/:{id}", RetreivePackage)
	}
	
	packageList := router.Group("/packages")
	{
		packageList.POST("", GetPackageList)	
	}

	//api.GET("", CreatePackage)
	router.Run(":8000")
}

//USEFUL FOR THINGS LIKE PAGE OFFSET IE localhost:8000/package?=1
/*func QueryParams(c *gin.Context) {
	queryPairs := c.Request.URL.Query()
	for key, values := range queryPairs {
		fmt.Printf("key: %v, val %v\n", key, values)
	}
}*/

func RetreivePackage(c *gin.Context){
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H {
		//access {id} from dynamic route which can be passed into db for processing
		"param": c.Param("{id}"),
	})
	//errors 400 (missing PackageID/Auth token/not authorized) or 404 (package not found)
}

func GetPackageList(c *gin.Context) {
	//array is needed in data -> how does this work?!

	//struct for json
	var ex JSON_LIST

	//set headers to application/json and require authentication
	c.Header("Content-Type", "application/json")
	//c.Request.Header.Add("X-Authorization", auth_token)

	//take data from frontend input and bind to json struct
	c.BindJSON(&ex)
	c.JSON(200, gin.H{"data": []interface{}{ex}})
	
	//add error codes 400 and 413

}

func CreatePackage(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H {
		"message": "create package has not been implemented.",
	})
}

func Reset(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H {
		"message": "delete package has not been implemented.",
	})
}