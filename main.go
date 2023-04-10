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

type GetPackageListJSON struct {
	Version string `json:"Version"`
	Name string `json:"Name"`
}

type PackageCreate struct {
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

/*func DBInit() {
	psqlinfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)

	db, err := sql.Open("postgres", psqlinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("db connected :)")
}*/

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
		//api.
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
	GetZip(url.URL)
	b64_string := EncodeZipFile()
	split := strings.Split(url.URL, "/")
	repo := split[len(split)-1]

	newObject := models.PackageCreate{Name: repo, Content: b64_string}
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
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H {
		"message": "delete package has not been implemented.",
	})
}