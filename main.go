package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/skoslitz/lid-backend/lidlib"
	"github.com/spf13/viper"
)

var applicationRoot string
var contentRoot string
var adminRoot string

func init() {

	// set lid repo root path
	err := os.Chdir("/home/kossi/lid-site/")

	if err != nil {
		fmt.Println("Can't change working directory")
	}

	applicationRoot, _ = os.Getwd()
	contentRoot = applicationRoot + "/content/"
	adminRoot = "/home/kossi/lid-frontend/app"

	// TODO:
	// case contentRoot = applicationRoot + "//content/"
	// should be false

	// check if content path is valid
	if _, err := os.Stat(contentRoot); os.IsNotExist(err) {
		fmt.Println("Content path is not valid. Please check!")
	}

}

func main() {

	// setup config file
	viper.SetConfigName("config")
	viper.ReadInConfig()

	// set config defaults
	viper.SetDefault("ContentDir", contentRoot)
	viper.SetDefault("AdminDir", adminRoot)
	viper.SetDefault("AssetsDir", applicationRoot+"/static")

	contentDir := viper.GetString("ContentDir")
	assetsDir := viper.GetString("AssetsDir")

	// create router
	router := NewRouter(&RouterConfig{
		Handlers: &Handlers{
			Config:     lidlib.NewConfig(applicationRoot + "/config.toml"),
			Dir:        lidlib.NewDir(),
			Page:       lidlib.NewPage(),
			ContentDir: contentDir,
			AssetsDir:  assetsDir,
		},
		AdminDir: viper.GetString("AdminDir"),
	})

	// start http server
	fmt.Println("Starting server on localhost:1313")
	fmt.Println("Content in ", contentDir)
	fmt.Println("Assets in ", assetsDir)
	fmt.Println("Admin in ", adminRoot)
	log.Fatal(http.ListenAndServe("localhost:1313", router))

}
