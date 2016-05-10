package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var applicationRoot string
var contentRoot string

type Content []string

func init() {

	// set lid repo root path
	err := os.Chdir("/home/kossi/lid-site/")

	if err != nil {
		fmt.Println("Can't change working directory")
	}

	applicationRoot, _ = os.Getwd()
	contentRoot = applicationRoot + "/content"

}

func readContentType(contentType string) {
	files, err := ioutil.ReadDir(contentRoot + contentType)

	if err != nil {
		log.Fatalln("Failed to open:", err)
	}

	var contentFile Content
	for _, file := range files {
		contentFile = append(contentFile, strings.TrimSuffix(file.Name(), ".md"))
	}

	data, _ := json.MarshalIndent(contentFile, "", "  ")

	fmt.Printf("%s\n", data)

}

func readConfig() {
	file, err := ioutil.ReadFile(applicationRoot + "/config.toml")

	if err != nil {
		log.Fatalln("Failed to open:", err)
	}
	fmt.Println(string(file))
}

func main() {
	readContentType("/")
	//readConfig()
}
