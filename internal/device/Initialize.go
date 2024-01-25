package device

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func (d *Device) Initialize() error {
	fmt.Println("Enter wifi ssid (then press enter):")

	// collect ssid and password from stdio
	reader := bufio.NewReader(os.Stdin)

	ssidInput, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	fmt.Println("\nEnter wifi password (then press enter):")
	pskInput, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// trim newline characters from user input
	ssidInput = strings.TrimSuffix(ssidInput, "\n")
	pskInput = strings.TrimSuffix(pskInput, "\n")

	success := d.TestWifiConfiguration(ssidInput, pskInput)
	if !success {
		return errors.New("Unable to connect to wifi with the provided credentials")
	}

	fmt.Println("Connection successful!")
	return nil
}
