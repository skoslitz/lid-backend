package lidlib

import (
	"fmt"
	"os/exec"
	"reflect"
	"strings"
)

func RunHugo(repoPath string, webFolder string) ([]byte, error) {
	hugo := exec.Command("hugo", "--source="+repoPath, "--destination="+webFolder)

	output, err := hugo.Output()
	if err != nil {
		return nil, err
	}

	return output, nil
}

func RunHugoPreview(baseUrlPrefix string, repoPath string) (string, error) {
	baseURL := strings.Join([]string{baseUrlPrefix, "/preview/"}, "")
	hugo := exec.Command("hugo", "--source="+repoPath, "--baseURL="+baseURL, "--canonifyURLs=true")

	output, err := hugo.Output()
	if err != nil {
		fmt.Println("Cmd hugo error output: ")
		fmt.Println(hugo.Stdout)
		fmt.Println(reflect.TypeOf(hugo.Stdout))
		fmt.Println("**********")
		output := fmt.Sprint(hugo.Stdout)
		fmt.Println(reflect.TypeOf(output))
		fmt.Println("**********")
		fmt.Println("hugo.Stderr")
		fmt.Println(hugo.Stderr)
		// make hugo.Stdout *bytes.Buffer to string and return
		return output, err
	}

	_output := fmt.Sprint(output)
	return _output, nil
}
