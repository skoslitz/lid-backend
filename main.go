package main

import (
	"log"
	"net/http"
	"os"
	/*	"os/exec"
		"runtime"*/

	"github.com/fatih/color"
	"github.com/skoslitz/lid-backend/lidlib"
	"github.com/spf13/viper"
	"github.com/stevedomin/termtable"
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
		messenger := color.New(color.Bold, color.FgRed).PrintlnFunc()
		messenger("LiD-online Repo konnte nicht gefunden werden. Bitte die config.toml prüfen!")
		log.Fatal()
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

	// lid-backend dashboard
	t := termtable.NewTable(nil, &termtable.TableOptions{
		Padding:      1,
		UseSeparator: true,
	})
	// t.SetHeader([]string{"LOWERCASE", "", ""})
	t.AddRow([]string{"lid-repo/content", contentDir})
	t.AddRow([]string{"lid-repo/static", assetsDir})
	t.AddRow([]string{"lid-repo/public", previewDir})
	t.AddRow([]string{"lid-frontend", adminDir})
	messenger := color.New(color.Bold, color.FgGreen).PrintlnFunc()
	messenger("+--------------------------------------------------------------+")
	messenger("              LiD Online Content API                            ")
	messenger(t.Render())
	messenger("  Serveradresse: localhost:1313                                 ")
	messenger("  Browse http://localhost:1313/api/dir/regionen                        ")
	messenger("+--------------------------------------------------------------+")

	// start http server
	log.Fatal(http.ListenAndServe("localhost:1313", router))

}
