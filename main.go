package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var applicationRoot string
var contentRoot string

func init() {

	// set lid repo root path
	err := os.Chdir("/home/kossi/lid-site/")

	if err != nil {
		fmt.Println("Can't change working directory")
	}

	applicationRoot, _ = os.Getwd()
	contentRoot = applicationRoot + "/content/"

}

func main() {

	router := NewRouter()

	log.Fatal(http.ListenAndServe("localhost:1313", router))

}
