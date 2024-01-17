package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	if isProdEnvironment() {
		fmt.Printf("Prod environment detected\n\n")
	} else {
		fmt.Printf("Dev environment detected\n\n")
	}
}

func isProdEnvironment() bool {
	// detect the linux kernel; if it doesn't include the text "raspberrypi" then it's a test device / environment
	byteOutput, err := exec.Command("uname", "-a").Output()

	if err != nil {
		log.Fatal(err)
		return false
	}

	kernelDescription := string(byteOutput[:])

	if strings.Contains(kernelDescription, "raspberrypi") {
		return true
	}

	return false
}
