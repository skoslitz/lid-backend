package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

func readContentType(contentType string) {
	files, err := ioutil.ReadDir(apiRoot + contentType)

	if err != nil {
		log.Fatalln("Failed to open:", err)
	}

	var content []string
	for _, file := range files {
		content = append(content, strings.TrimSuffix(file.Name(), ".md"))
	}

	data, _ := json.MarshalIndent(content, "", "  ")

	fmt.Printf("%s\n", data)

}

func main() {
	readContentType("/exkursionen")
}
