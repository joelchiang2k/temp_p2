package HelperFunctions

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

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

}