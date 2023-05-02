package main

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"ex/part2/HelperFunctions"
	"ex/part2/controllers"
	"ex/part2/logger"
	"ex/part2/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"

	//"database/gcpbucket"

	"github.com/gin-gonic/gin"
)

//project id: perceptive-tape-383118

/*type PackageMetadata struct {
	//Version string `json:"Version"`
	Name string `json:"Name"`
	//ID string `json:"id"`
}*/

type PackageJsonInfo struct {
	Homepage string `json:"homepage"`
	Version  string `json:"Version"`
}

type PackageCreate struct {
	Name    string `json:"packageName"`
	Version string `json:"packageVersion"`
	Content string `json:"Content"`
	URL     string `json:"URL"`
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

	logger := logger.GetInst()
	logger.Println("program beginning")

	api := router.Group("/package")
	{
		api.POST("", CreatePackage)
		api.GET("/:{id}", RetreivePackage)
		api.PUT("/:{id}", controllers.UpdatePackage)
		api.DELETE("/:{id}", DeletePackageById)
		//api.GET("/:{id}/rate, RatePackage")
		api.POST("/byRegEx", controllers.ByRegex)
	}

	packageList := router.Group("/packages")
	{
		packageList.POST("", GetPackageList)
	}

	resetRoute := router.Group("/reset")
	{
		resetRoute.DELETE("", Reset)
	}

	auth := router.Group("/authenticate")
	{
		auth.PUT("", controllers.Authenticate)
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

func RetreivePackage(c *gin.Context) {
	//c.Header("Content-Type", "application/json")

	//create variable to hold response data
	var packageToRetreive models.PackageCreate

	//change so that if id is missing return error
	if c.Param("{id}") == "/" {
		c.JSON(400, "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
	} else if err := models.DB.Where("id = ?", c.Param("{id}")).First(&packageToRetreive).Error; err != nil {
		c.JSON(404, "Package does not exist.")
	} else {
		//logging purposes
		niceJSON, err := json.MarshalIndent(packageToRetreive, "", " ")
		if err != nil {
			fmt.Println(err)
			return

		}	
		fmt.Println("/package/:{id} GET response")

		fmt.Println(string(niceJSON))
		c.JSON(200, gin.H{
			"metadata": gin.H{
				"Name":    packageToRetreive.Name,
				"Version": packageToRetreive.Version,
				"ID":      packageToRetreive.ID,
			},
			"data": gin.H{
				"Content": packageToRetreive.Content,
			},
		})
		return
	}
	//errors 400 (missing PackageID/Auth token/not authorized) or 404 (package not found)
}

func GetPackageList(c *gin.Context) {
	//array is needed in data -> how does this work?!

	//struct for json
	var ex PackageCreate

	//set headers to application/json and require authentication
	//c.Header("Content-Type", "application/json")
	//c.Request.Header.Add("X-Authorization", auth_token)

	//take data from frontend input and bind to json struct
	c.BindJSON(&ex)

	//c.JSON(200, gin.H{"data": []interface{}{ex}})
	c.JSON(500,"not implemented yet")

	//add error codes 400 and 413
}

func RatePackage(c *gin.Context) {
	var packageToRate models.PackageCreate

	if c.Param("{id}") == "/" {
		c.JSON(400, "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.")
	} else if err := models.DB.Where("id = ?", c.Param("{id}")).First(&packageToRate).Error; err != nil {
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

func CreatePackage(c *gin.Context) {
	//creates new variable of {struct} type and binds data from incoming request to new variable
	//returns error on bad req
	logger := logger.GetInst()
	var newPackage PackageCreate

	if err := c.BindJSON(&newPackage); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	niceJSON, err := json.MarshalIndent(newPackage, "", " ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("package POST request")
	fmt.Println(string(niceJSON))

	logger.Printf("Incoming Request for /package POST \nContent: %s\nURL: %s\n", newPackage.Content, newPackage.URL)

	if newPackage.URL != "" && newPackage.Content != "" {
		c.JSON(400, "URL and Content both set")
	} else if newPackage.URL == "" && newPackage.Content == "" {
		c.JSON(400, "Neither URL or Content set.")
	} else if newPackage.URL != "" {
		//process zip and upload to db

		HelperFunctions.GetZip(newPackage.URL)

		//get info from package.json
		var packageJsonObj PackageJsonInfo
		getPackageJsonInfo(&packageJsonObj)

		b64_string := HelperFunctions.EncodeZipFile()
		split := strings.Split(newPackage.URL, "/")
		repo := split[len(split)-1]

		newObject := models.PackageCreate{Name: repo, Version: packageJsonObj.Version, Content: b64_string, URL: newPackage.URL}

		//check if package already exists
		if result := models.DB.Where("name = ?", newObject.Name).First(&newObject).RowsAffected; result == 1 {
			c.JSON(409, "Package exists already.")	
			return
		}

		models.DB.Create(&newObject)

		//newPackage only used for incoming request -> GET ID FROM newObject

		//print response
		niceJSON2, err := json.MarshalIndent(newObject, "", " ")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("package POST INGEST response:")
		fmt.Println(string(niceJSON2))

		//attempt to log
		logger.Printf("Package Ingest Response: \n metadata:\n	Name: %s\n	Version: %s\n	ID: %d\n data:\n	Content: %s\n", newObject.Name, newObject.Version, newObject.ID, newObject.Content)

		//response w/ 201 and attributes
		c.JSON(201, gin.H{
			"metadata": gin.H{
				"Name":    newObject.Name,
				"Version": newObject.Version,
				"ID":      newObject.ID,
			},
			"data": gin.H{
				"Content": newObject.Content,
			},
		})
	} else if newPackage.Content != "" {
		decodedString, err := base64.StdEncoding.DecodeString(newPackage.Content)

		if err != nil {
			panic(err)
		}

		directory, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}

		//create out file in zip_files dir
		outFile, err := os.Create(directory + "/zip_files/out.zip")
		if err != nil {
			fmt.Println(err)
		}
		defer outFile.Close()

		_, err = outFile.Write(decodedString)
		if err != nil {
			fmt.Println(err)
		}

		var packageJsonObj PackageJsonInfo
		getPackageJsonInfo(&packageJsonObj)

		split := strings.Split(packageJsonObj.Homepage, "/")
		repo := split[len(split)-1]

		newObject := models.PackageCreate{Name: repo, Version: packageJsonObj.Version, Content: newPackage.Content, URL: packageJsonObj.Homepage}

		//check if package already exists
		if result := models.DB.Where("name = ?", newObject.Name).First(&newObject).RowsAffected; result == 1 {
			c.JSON(409, "Package exists already.")	
			return
		}

		models.DB.Create(&newObject)

		niceJSON2, err := json.MarshalIndent(newObject, "", " ")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("package POST Upload by Zip response:")
		fmt.Println(string(niceJSON2))

		c.JSON(201, gin.H{
			"metadata": gin.H{
				"Name":    newObject.Name,
				"Version": newObject.Version,
				"ID":      newObject.ID,
			},
			"data": gin.H{
				"Content": newObject.Content,
			},
		})

	}

}

/*func EncodeZipFile() (b64 string) {
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
}*/

/*func GetZip(url string) {

	//split owner and repository strings from original url for later use
	split := strings.Split(url, "/")
	owner := split[len(split)-2]
	repo := split[len(split)-1]

	//get current working directory for use in filepath
	directory, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	//create out file in zip_files dir
	outFile, err := os.Create(directory + "/zip_files/out.zip")
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()

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

}*/

func splitPaths(path string) []string {
	paths := strings.Split(strings.Replace(path, "\\", "/", -1), "/")

	if len(paths) > 1 && !strings.Contains(paths[len(paths)-1], ".") {
		paths = paths[:len(paths)-1]
	}

	return paths
}

func getPackageJsonInfo(packageJsonObj *PackageJsonInfo) {
	directory, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	//open zipfile
	zipFile, err := zip.OpenReader(directory + "/zip_files/out.zip")
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	//find name of root directory for use in unmarshalling package.json data
	directoryCounter := make(map[string]int)
	for _, fileName := range zipFile.File {
		directories := splitPaths(fileName.Name)
		if len(directories) > 0 {
			directoryCounter[directories[0]]++
		}
	}

	max := 0
	root := ""
	for i, count := range directoryCounter {
		if count > max {
			max = count
			root = i
		}
	}

	//fmt.Println(root)

	var data []byte
	for _, file := range zipFile.File {
		//fmt.Println(file.Name)
		if file.Name == (root + "/package.json") {
			f, err := file.Open()
			if err != nil {
				panic(err)
			}
			defer f.Close()
			data, err = ioutil.ReadAll(f)
			if err != nil {
				panic(err)
			}
			break
		}
	}

	if data == nil {
		panic("No package.json found in github repository")
	}

	err = json.Unmarshal(data, packageJsonObj)
	if err != nil {
		panic(err)
	}
}

func Reset(c *gin.Context) {

	if tx := models.DB.Exec("TRUNCATE TABLE package_creates RESTART IDENTITY"); tx.Error != nil {
		//change
		panic(tx.Error)
	}

	c.JSON(200, "Registry is reset.")
}
