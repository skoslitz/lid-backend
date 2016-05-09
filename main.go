package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var apiRoot string

func init() {

	// set api root path
	err := os.Chdir("/home/kossi/lid-site/content")

	if err != nil {
		fmt.Println("Can't change working directory")
	}

	apiRoot, _ = os.Getwd()

}

func main() {

	files, err := ioutil.ReadDir(apiRoot + "/themen")

	if err != nil {
		log.Fatalln("Failed to open:", err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}
