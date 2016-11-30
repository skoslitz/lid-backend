package lidlib

import (
	"os/exec"
	"strings"
)

func RunHugo(repoPath string, webFolder string) ([]byte, error) {
	hugo := exec.Command("hugo17", "--source="+repoPath, "--destination="+webFolder)

	output, err := hugo.Output()
	if err != nil {
		return output, nil
	}

	return output, nil
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
