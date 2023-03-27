package gotest

import (
	"io/ioutil"
	"log"
	"net/http"
)

func gotest() {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
	   log.Fatalln(err)
	}
 //We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}
 }