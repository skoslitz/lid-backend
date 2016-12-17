package lidlib

import (
	//"golang.org/x/sys/unix"
	//"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunHugo(repoPath string, webFolder string) ([]byte, error) {

	// TODO: make this working on sys/windows
	//permission := unix.Access(webFolder, unix.W_OK) == nil

	//if permission {

	//change to repoPath with existing hugo17 executable
	os.Chdir(repoPath)
	hugo := exec.Command("hugo17", "--source="+repoPath, "--destination="+webFolder)
	output, err := hugo.Output()
	if err != nil {
		return output, nil
	}

	return output, nil
	/*} else {
		message := "Keine Berechtigung!"
		return []byte(message), nil
	}*/

}

func RunHugoPreview(repoPath string, baseUrlPrefix string) ([]byte, error) {
	baseURL := strings.Join([]string{baseUrlPrefix, "/preview/"}, "")
	hugo := exec.Command("hugo17", "--source="+repoPath, "--baseURL="+baseURL, "--canonifyURLs=true")

	output, err := hugo.Output()
	if err != nil {
		return output, nil
	}

	return output, nil
}
