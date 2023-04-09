package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	//"database/gcpbucket"
	"cloud.google.com/go/storage"
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

/*func DbInit(c gin.Context){
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println(err)
	}
	c.Set("gcpBucket", client)
	c.Next()
}*/

func CORS(c *gin.Context) {

	// First, we add the headers with need to enable CORS
	// Make sure to adjust these headers to your needs
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	/*ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println(err)
	}
	c.Set("gcpBucket", client)*/

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

	/*ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println(err)
	}
	client = c*/

	//router.Use(static.Serve("/", static.LocalFile("./src", true)))
	router.Use(CORS)
	//router.Use(DbInit)
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
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println(err)
	}
	//ctx, cancel := context.WithTimeout(ctx, time.Second*10)	
	//defer cancel()

	//creates new variable of {struct} type and binds data from incoming request to new variable
	//returns error on bad req
	var url PackageCreate

	if err := c.BindJSON(&url); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	//process zip and upload to db
	GetZip(url.URL)
	//b64_string := EncodeZipFile()

	var projectID string = "perceptive-tape-383118"

	//client, ok:= c.MustGet("gcpBucket").(*storage.Client)
	/*if !ok {
		fmt.Println(ok)
	}*/

	bkt := client.Bucket("new-461-bucket")
	if err := bkt.Create(ctx, projectID, nil); err != nil {
		// TODO: Handle error.
		fmt.Println(err)
	}
	/*obj := bucket.Object("data")
	w := obj.NewWriter(ctx)
	if _, err := fmt.Fprintf(w, "This object contains text.\n"); err != nil {
		// TODO: Handle error.
		fmt.Println(err)
	}
	if err := w.Close(); err != nil {
		fmt.Println(err)
	}*/

	//response
	c.JSON(200, gin.H{
		"url": url.URL,
	})
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