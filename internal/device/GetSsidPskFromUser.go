package device

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (d *Device) GetSsidPskFromUser() (string, string, error) {
	fmt.Println("Enter wifi ssid (then press enter):")

	// collect ssid and password from stdio
	reader := bufio.NewReader(os.Stdin)

	ssidInput, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Println("\nEnter wifi password (then press enter):")
	pskInput, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	// trim newline characters from user input
	ssidInput = strings.TrimSuffix(ssidInput, "\n")
	pskInput = strings.TrimSuffix(pskInput, "\n")

	return ssidInput, pskInput, nil
}
