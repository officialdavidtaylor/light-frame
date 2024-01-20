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
	InitialSetupCompleted bool `json:"initialSetupCompleted"`
}

func main() {
	// initialize device and instantiate depending on environment
	var d device.D

	if isProdEnvironment() {
		fmt.Printf("Prod environment detected\n\n")
		d = device.NewProdDevice(5)
	} else {
		fmt.Printf("Dev environment detected\n\n")
		d = device.NewDevDevice(5)
	}

	config := loadConfiguration()

	if !config.InitialSetupCompleted {
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
	// detect cwd
	p, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to determine working directory")
	}
	// open / create config file
	configFileContents, err := os.ReadFile(p + "/cmd/conf.json")
	// if unable to open the configuration, try to create the file
	if err != nil {
		fmt.Println(err)

		f, e := os.Create(p + "/cmd/conf.json")
		if e != nil {
			log.Fatal("Could not create configuration file")
		}
		// after creating the file, close it and return an empty configuration struct
		f.Close()

		return Configuration{}
	}

	// parse JSON
	configuration := Configuration{}
	unmarshalErr := json.Unmarshal(configFileContents, &configuration)

	if unmarshalErr != nil {
		log.Fatal("Failure parsing configuration file.")
	}

	return configuration
}
