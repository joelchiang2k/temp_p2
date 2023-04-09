package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetPackageListJSON struct {
	Version string `json:"Version"`
	Name string `json:"Name"`
}

type PackageCreate struct {
	Content string `json:"Content`
	URL string `json:"URL"`
	//JSProgram
}

func CORS(c *gin.Context) {

	// First, we add the headers with need to enable CORS
	// Make sure to adjust these headers to your needs
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	// Second, we handle the OPTIONS problem
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		//If options request is received, return 200 so that the program can process
		// a different request
		c.AbortWithStatus(http.StatusOK)
	}
}

func main() {
	router := gin.Default()

	//router.Use(static.Serve("/", static.LocalFile("./src", true)))
	router.Use(CORS)
	api := router.Group("/package")
	{
		api.POST("", CreatePackage)
		api.GET("/:{id}", RetreivePackage)
		//api.PUT("/:{id}", ADD FUNC FOR PUT)
		//api.DELETE("/:{id}", ADD FUN FOR DELETE)
		//api.
	}
	
	packageList := router.Group("/packages")
	{
		packageList.POST("", GetPackageList)	
	}

	resetRoute := router.Group("/reset")
	{
		reset.DELETE("", Reset)
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
	var ex GetPackageListJSON

	//set headers to application/json and require authentication
	c.Header("Content-Type", "application/json")
	//c.Request.Header.Add("X-Authorization", auth_token)

	//take data from frontend input and bind to json struct
	c.BindJSON(&ex)
	c.JSON(200, gin.H{"data": []interface{}{ex}})
	
	//add error codes 400 and 413
}

func CreatePackage(c *gin.Context) {
	
	//creates new variable of {struct} type and binds data from incoming request to new variable
	//returns error on bad req
	var url PackageCreate

	if err := c.BindJSON(&url); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	//process zip and upload to db

	//response
	c.JSON(200, gin.H{
		"url": url.URL,
	})
}

func Reset(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H {
		"message": "delete package has not been implemented.",
	})
}