package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSON_LIST struct {
	Version string `json:"Version"`
	Name string `json:"Name"`
}

type PackageCreate struct {
	URL string `json:"url"`
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
        
		// Everytime we receive an OPTIONS request, 
		// we just return an HTTP 200 Status Code
		// Like this, Angular can now do the real 
		// request using any other method than OPTIONS
		c.AbortWithStatus(http.StatusOK)
	}
}

func main() {
	router := gin.Default()

	//router.Use(static.Serve("/", static.LocalFile("./src", true)))
	router.Use(CORS)
	api := router.Group("/package")
	{
		/*api.GET("/", func(c *gin.Context){
			c.JSON(http.StatusOK, gin.H{
				"message": "hello",
			})
		})*/
		api.POST("", CreatePackage)
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
	
	//creates new variable of {struct} type and binds data from incoming request to new variable
	//returns error on bad req
	var url PackageCreate

	if err := c.BindJSON(&url); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

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