package lidlib

import (
	"fmt"
	"os/exec"
	"strings"
)

func RunHugo(repoPath string, webFolder string) ([]byte, error) {
	hugo := exec.Command("hugo", "--source="+repoPath, "--destination="+webFolder)

	output, err := hugo.Output()
	if err != nil {
		return []byte(fmt.Sprint(hugo.Stdout)), err
	}

	return output, nil
}

func RunHugoPreview(baseUrlPrefix string, repoPath string) ([]byte, error) {
	baseURL := strings.Join([]string{baseUrlPrefix, "/preview/"}, "")
	hugo := exec.Command("hugo", "--source="+repoPath, "--baseURL="+baseURL, "--canonifyURLs=true")

	output, err := hugo.Output()
	if err != nil {
		return []byte(fmt.Sprint(hugo.Stdout)), err
	}

	//_output := fmt.Sprint(output)
	return output, nil
}
