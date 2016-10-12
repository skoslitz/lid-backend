package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/skoslitz/lid-backend/lidlib"
	"github.com/spf13/viper"
)

func init() {

	configRoot, _ := os.Getwd()

	// setup config file
	// Find and read the config
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.AddConfigPath(configRoot) // path to look for the config file in
	viper.ReadInConfig()            // Find and read the config file

}

func main() {

	applicationRoot := viper.GetString("repopath")
	contentDir := viper.GetString("contentpath")
	assetsDir := viper.GetString("assetspath")
	previewDir := viper.GetString("previewpath")
	//adminDir := viper.GetString("adminpath")

	// check if content path is valid
	if _, err := os.Stat(contentDir); os.IsNotExist(err) {
		fmt.Println("LiD Inhaltspfad konnte nicht gefunden werden. Bitte die config.toml prüfen!")
	}

	// create router
	router := NewRouter(&RouterConfig{
		Handlers: &Handlers{
			Config:     lidlib.NewConfig(applicationRoot + "/config.toml"),
			Dir:        lidlib.NewDir(),
			Page:       lidlib.NewPage(),
			ContentDir: contentDir,
			AssetsDir:  assetsDir,
			PreviewDir: previewDir,
		},
		AdminDir: viper.GetString("adminpath"),
	})

	// start http server
	fmt.Println("Server gestartet auf localhost:1313")
	fmt.Println("LiD Inhaltspfad: ", contentDir)
	fmt.Println("LiD Anhangspfad: ", assetsDir)
	fmt.Println("LiD Vorschaupfad ", previewDir)
	//fmt.Println("Admin in ", adminDir)
	log.Fatal(http.ListenAndServe("localhost:1313", router))

}
