package lidlib

import (
	"fmt"
	"os/exec"
	"strings"
)

func RunHugo(repoPath string, webFolder string) (string, error) {
	hugo := exec.Command("hugo", "--source="+repoPath, "--destination="+webFolder)

	_, err := hugo.Output()
	if err != nil {
		//return []byte(fmt.Sprint(hugo.Stdout)), err
		return fmt.Sprint(hugo.Stdout), nil
	}

	//return output, nil
	return fmt.Sprint(hugo.Stdout), nil
}

func RunHugoPreview(baseUrlPrefix string, repoPath string) (string, error) {
	baseURL := strings.Join([]string{baseUrlPrefix, "/preview/"}, "")
	hugo := exec.Command("hugo17", "--source="+repoPath, "--baseURL="+baseURL, "--canonifyURLs=true")

	_, err := hugo.Output()
	if err != nil {
		//return []byte(fmt.Sprint(hugo.Stdout)), err
		return fmt.Sprint(hugo.Stdout), nil
	}

	//return output, nil
	return fmt.Sprint(hugo.Stdout), nil
}
