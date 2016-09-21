package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

func shortcodeFileName(f string) string {
	return strings.Split(f, "-")[1]
}

func containsContentType(t string) bool {
	ctypes := [...]string{"themen", "exkursionen"}
	for _, v := range ctypes {
		if v == t {
			return true
		}
	}
	return false
}

func fileExists(fp string) bool {
	info, err := os.Stat(fp)
	if err != nil {
		return false
	}
	return info.IsDir() == false
}

func dirExists(fp string) bool {
	info, err := os.Stat(fp)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func printError(w http.ResponseWriter, err interface{}) {
	printJson(w, err)
}

func printJson(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/vnd.api+json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	json.NewEncoder(w).Encode(obj)
}
