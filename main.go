package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"ex/part2/models"

	_ "github.com/lib/pq"

	//"database/gcpbucket"

	"github.com/gin-gonic/gin"
)

//project id: perceptive-tape-383118

type PackageMetadata struct {
	//Version string `json:"Version"`
	Name string `json:"Name"`
	//ID string `json:"id"`
}

type PackageCreate struct {
	Name string `json:"packageName"`
	Version string `json:"packageVersion"`
	Content string `json:"Content"`
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

	//DBInit()
	models.ConnectDatabase()
	router.Use(CORS)
	api := router.Group("/package")
	{
		api.POST("", CreatePackage)
		api.GET("/:{id}", RetreivePackage)
		//api.PUT("/:{id}", ADD FUNC FOR PUT)
		api.DELETE("/:{id}", DeletePackageById)
		//api.GET("/:{id}/rate, RatePackage")
	}
	
	packageList := router.Group("/packages")
	{
		packageList.POST("", GetPackageList)	
	}

	resetRoute := router.Group("/reset")
	{
		resetRoute.DELETE("", Reset)
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

func DeletePackageById(c *gin.Context) {
	var packageToDelete models.PackageCreate
	
	if err := models.DB.Where("id = ?", c.Param("{id}")).First(&packageToDelete).Error; err != nil {
		c.JSON(404, "Package does not exist.")
	}

	models.DB.Delete(&packageToDelete)
	c.JSON(200, "Package is deleted.")
}

func RetreivePackage(c *gin.Context){
	//c.Header("Content-Type", "application/json")

	//create variable to hold response data
	var packageToRetreive models.PackageCreate
	
	//change so that if id is missing return error
	if c.Param("{id}") == "/"{
		c.JSON(400, "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
	}else if err := models.DB.Where("id = ?", c.Param("{id}")).First(&packageToRetreive).Error; err != nil {
		c.JSON(404, "Package does not exist.")
	}else{ 
		c.JSON(200, gin.H {
			//access {id} from dynamic route which can be passed into db for processing
			"data": []interface{}{packageToRetreive},
		})
	}
	//errors 400 (missing PackageID/Auth token/not authorized) or 404 (package not found)
}

func GetPackageList(c *gin.Context) {
	//array is needed in data -> how does this work?!

	//struct for json
	var ex PackageMetadata

	//set headers to application/json and require authentication
	c.Header("Content-Type", "application/json")
	//c.Request.Header.Add("X-Authorization", auth_token)

	//take data from frontend input and bind to json struct
	c.BindJSON(&ex)
	c.JSON(200, gin.H{"data": []interface{}{ex}})
	
	//add error codes 400 and 413
}

func RatePackage(c *gin.Context) {
	var packageToRate models.PackageCreate	

	if c.Param("{id}") == "/"{
		c.JSON(400, "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
	}else if err := models.DB.Where("id = ?", c.Param("{id}")).First(&packageToRate).Error; err != nil {
		c.JSON(404, "Package does not exist.")
	}

	//DANIEL ENTER RATE STUFF HERE
	//could insert functions under a subdirectory called rateFunctions/
	//then call like rateFunctions.netScore(), rateFunctions.Responsiveness(), etc...

	//if anything in rating funcs fail
	//c.JSON(500, "The package rating system choked on at least one of the metrics.")

	//else everything ok
	//c.JSON(200, gin.H{
		//ratingStruct
	//})
}

//finish
func UpdatePackage(c *gin.Context) {
	var pkg models.PackageCreate

	//get package if exists in db
	//incorrect format in route string
	if c.Param("{id}") == "/"{
		c.JSON(400, "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
	}else if err := models.DB.Where("id = ?", c.Param("{id}")).First(&pkg).Error; err != nil {
		//package not found -> return error 404
		c.JSON(404, "Package does not exist.")
	}

	var packageToUpdate PackageCreate
	//validate name, ID, version are matched
	if err := c.BindJSON(&packageToUpdate); err != nil {
		c.AbortWithError(400, err)
		return
	}
	//ADD VERSION
	//if(packageToUpdate.Name == pkg.Name && c.Param("{id}") == string(pkg.ID)){
		
	//models.DB.Update("Content", packageToUpdate.Content)
	//}

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
	GetZip(url.URL)
	b64_string := EncodeZipFile()
	//split := strings.Split(url.URL, "/")
	//repo := split[len(split)-1]

	newObject := models.PackageCreate{Name: url.Name, Version: url.Version, Content: b64_string, URL: url.URL}
	models.DB.Create(&newObject)

	//response
	/*c.JSON(200, gin.H{
		"url": url.URL,
	})*/
	c.JSON(201, gin.H{"data": newObject})
	
}

func EncodeZipFile() (b64 string) {
	//get working directory
	directory, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	file, _ := os.Open(directory + "/zip_files/out.zip")

	reader := bufio.NewReader(file)
	content, _ := ioutil.ReadAll(reader)

	encodedString := base64.StdEncoding.EncodeToString(content)

	return encodedString
}


func GetZip(url string) {

	//split owner and repository strings from original url for later use
	split := strings.Split(url, "/")
	owner := split[len(split)-2]
	repo := split[len(split)-1]

	//get current working directory for use in filepath
	directory, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(directory)

	//create out file in zip_files dir
	outFile, err := os.Create(directory + "/zip_files/out.zip")
	if err != nil {
		fmt.Println(err)
	}

	//send get request to correct git repo and get zip file contents
	resp, err := http.Get("https://github.com/" + owner + "/" + repo + "/archive/master.zip")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()	

	//write response body to zip file created above
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	
}

func Reset(c *gin.Context) {
	
	if tx := models.DB.Exec("TRUNCATE TABLE package_creates RESTART IDENTITY"); tx.Error != nil {
		//change
		panic(tx.Error)
	}

	c.JSON(200, "Registry is reset.")
}