package sys

import (
	"os/exec"
	"strings"
)

func IsProdEnvironment() bool {
	// detect the linux kernel; if it doesn't include the text "raspberrypi" then it's a test device / environment
	byteOutput, err := exec.Command("uname", "-a").Output()
	if err != nil {
		return false
	}

	kernelDescription := string(byteOutput[:])

	if strings.Contains(kernelDescription, "raspberrypi") {
		return true
	}

	return false
}
