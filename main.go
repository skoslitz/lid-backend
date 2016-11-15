package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	/*	"os/exec"
		"runtime"*/

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
	adminDir := viper.GetString("adminpath")

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

	// open browser with predefined url
	/*switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", "http://localhost:1313/api/dir/regionen").Start()
	case "windows":
		exec.Command("cmd", "/c", "start", "http://localhost:1313/api/dir/regionen").Start()
	case "darwin":
		exec.Command("open", "http://localhost:1313/api/dir/regionen").Start()
	default:
		fmt.Errorf("unsupported platform")
	}*/

	// start http server
	fmt.Println("LiD Inhaltsschnittstelle wird gestartet. -- server: localhost:1313 --")
	fmt.Println("------------------------------")
	fmt.Println("LiD Inhaltspfad: ", contentDir)
	fmt.Println("LiD Anhangspfad: ", assetsDir)
	fmt.Println("LiD Vorschaupfad ", previewDir)
	//fmt.Println("Browser mit Regionenendpunkt wird geladen.")
	fmt.Println("Admin in ", adminDir)
	log.Fatal(http.ListenAndServe("localhost:1313", router))

}
