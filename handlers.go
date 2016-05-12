package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

type Content []string

func Index(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome!")
}

func ReadConfig(w http.ResponseWriter, r *http.Request) {

	file, err := ioutil.ReadFile(applicationRoot + "/config.toml")

	if err != nil {
		log.Fatalln("Failed to open:", err)
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(file))
}

func ReadContentIndex(w http.ResponseWriter, r *http.Request) {

	files, err := ioutil.ReadDir(contentRoot)

	if err != nil {
		log.Fatalln("Failed to open:", err)
	}

	var contentFile Content
	for _, file := range files {
		contentFile = append(contentFile, strings.TrimSuffix(file.Name(), ".md"))
	}

	data, _ := json.MarshalIndent(contentFile, "", "  ")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(data))

}

func ReadContentType(w http.ResponseWriter, r *http.Request) {

	contentType := mux.Vars(r)["contentType"]
	files, err := ioutil.ReadDir(contentRoot + contentType)

	if err != nil {
		log.Fatalln("Failed to open:", err)
	}

	var contentFile Content
	for _, file := range files {
		contentFile = append(contentFile, strings.TrimSuffix(file.Name(), ".md"))
	}

	data, _ := json.MarshalIndent(contentFile, "", "  ")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(data))

}

func ReadContentTypeFile(w http.ResponseWriter, r *http.Request) {

	contentType := mux.Vars(r)["contentType"]
	fileName := mux.Vars(r)["fileName"]
	filePath := []string{contentRoot, contentType, "/", fileName, ".md"}

	file, err := ioutil.ReadFile(strings.Join(filePath, ""))

	if err != nil {
		log.Fatalln("Failed to open:", err)
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(file))

}

func ReadSplitContentTypeFile(w http.ResponseWriter, r *http.Request) {

	contentType := mux.Vars(r)["contentType"]
	fileName := mux.Vars(r)["fileName"]
	filePart := mux.Vars(r)["filePart"]
	filePath := []string{contentRoot, contentType, "/", fileName, ".md"}

	file, err := ioutil.ReadFile(strings.Join(filePath, ""))

	if err != nil {
		log.Fatalln("Failed to open:", err)
	}

	searchTerm := `[+]{3}([\s\S]+?)[+]{3}([\s\S]*)`
	re := regexp.MustCompile(searchTerm)

	switch filePart {
	case "meta":
		fileMeta := re.FindStringSubmatch(string(file))[1]
		fileMetaJSON, _ := json.MarshalIndent(fileMeta, "", "  ")
		fmt.Fprintln(w, string(fileMetaJSON))
	case "content":
		fileContent := re.FindStringSubmatch(string(file))[2]
		fileContentJSON, _ := json.MarshalIndent(fileContent, "", "  ")
		fmt.Fprintln(w, string(fileContentJSON))
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "param does not exist")
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
