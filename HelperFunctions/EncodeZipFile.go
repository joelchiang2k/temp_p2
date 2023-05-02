package HelperFunctions

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

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