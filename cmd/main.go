package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"main/internal/device"
)

func main() {
	// initialize device and instantiate depending on environment
	var d device.D

	if isProdEnvironment() {
		fmt.Printf("Prod environment detected\n\n")
		d = &device.ProdDevice{}
	} else {
		fmt.Printf("Dev environment detected\n\n")
		d = &device.DevDevice{}
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
