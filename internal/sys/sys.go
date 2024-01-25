package sys

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
)

type C interface {
	Read()
	Write()
}

type Config struct {
	InitialSetupCompleted bool `json:"initialSetupCompleted"`

	SlideshowInterval int `json:"slideshowInterval"`
}

// Populate the config struct with values from the file saved on the disk
//
// If the file doesn't exist, create one with the current values in the struct
func (c *Config) Read() error {

	// detect cwd
	p, err := os.Getwd()
	if err != nil {
		return errors.New("Unable to determine working directory")
	}

	configFileContents, err := os.ReadFile(p + "/cmd/conf.json")
	// if unable to open the configuration, create the file
	if err != nil {
		return c.Write()
	}

	// parse JSON
	uErr := json.Unmarshal(configFileContents, &c)
	if uErr != nil {
		return errors.New("Failure parsing configuration file.")
	}

	return nil
}

// Marshal the current Config struct and save the file to the disk
func (c *Config) Write() error {
	co, coErr := json.Marshal(c)
	if coErr != nil {
		return errors.New("Failure marshaling new configuration struct into JSON byte[]")
	}

	wd, wdErr := os.Getwd()
	if wdErr != nil {
		return errors.New("Failure to determine working directory")
	}

	cfErr := os.WriteFile(wd+"/cmd/conf.json", co, 0666)
	if cfErr != nil {
		return errors.New("Could not write to config file.")
	}

	return nil
}

// Use the linux "clear" command to clear the terminal
func ClearScreen() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
