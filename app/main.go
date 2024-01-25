package main

import (
	"fmt"
	"log"
	"main/internal/device"
	"main/internal/sys"
)

func main() {
	sys.ClearScreen()

	// set default values that can be overridden with values saved in config file
	config := sys.Config{
		SlideshowInterval: 5,
	}

	rErr := config.Read()
	if rErr != nil {
		log.Fatal(rErr)
	}

	// initialize device and instantiate depending on environment
	var d device.D

	if sys.IsProdEnvironment() {
		d = device.NewProdDevice(config.SlideshowInterval)
	} else {
		fmt.Printf("Dev environment detected\n\n")
		d = device.NewDevDevice(config.SlideshowInterval)
	}

	for !config.InitialSetupCompleted {
		err := d.Initialize()
		if err != nil {
			continue
		}

		// update config value and save to disk
		config.InitialSetupCompleted = true
		err = config.Write()
		if err != nil {
			log.Fatal(err)
		}
	}

	d.StartSlideshow()
}
