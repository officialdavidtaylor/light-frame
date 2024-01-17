package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"main/internal/device"
	"main/internal/stateMachines"
)

type Configuration struct {
	initialSetupCompleted bool
}

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

	config := loadConfiguration()

	if !config.initialSetupCompleted {
		fmt.Println("Initial setup has _not_ been completed")

		// initialize state machine
		fmt.Printf("Initializing state machine...\n\n")

		fsm := stateMachines.InitWifiOnboardingWizard(d)

		err := fsm.Event(context.Background(), "next")

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("state at the end of the wifi config wizard: " + fsm.Current())
	}
		fmt.Println("Initial setup has been completed")
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

func loadConfiguration() Configuration {
	// open / create config file
	file, err := os.Open("conf.json")
	// if unable to open the configuration, try to create the file
	if err != nil {
		f, e := os.Create("conf.json")
		if e != nil {
			log.Fatal("could not create conf file")
		}
		fmt.Fprintf(f, "{}")
		// hoist the newly created file into the parent scope
		file = f
	}
	defer file.Close()
	// parse JSON
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	decodeErr := decoder.Decode(&configuration)

	if decodeErr != nil {
		fmt.Println("error:", decodeErr)
	}

	return configuration
}
