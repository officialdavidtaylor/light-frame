package sys

import (
	"encoding/json"
	"errors"
	"os"
)

type C interface {
	Read()
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
