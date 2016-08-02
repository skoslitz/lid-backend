package lidlib

import "os/exec"

func RunHugo() ([]byte, error) {
	hugo := exec.Command("hugo")

	output, err := hugo.Output()
	if err != nil {
		return nil, err
	}

	return output, nil
}

func RunHugoPreview() ([]byte, error) {
	baseURL := "http://localhost:1313/preview/"
	hugo := exec.Command("hugo", "--baseURL="+baseURL, "--canonifyURLs=true")

	output, err := hugo.Output()
	if err != nil {
		return nil, err
	}

	return output, nil
}
