package device

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

// Function for dev device / testing
func (d *Device) TestWifiConfiguration(ssid string, password string) bool {
	fmt.Println("TestWifiConfiguration called (via Device interface)")

	if ssid == "testSsid" && password == "testPassword" {
		return true
	}

	return false
}

// Use
func (d *ProdDevice) TestWifiConfiguration(ssid string, password string) bool {

	// duplicate wpa_supplicant file so we have a backup
	err := backupNetworkConfigFile()
	if err != nil {
		log.Fatal("Error backing up network config file.", err)
	}
	fmt.Println("wpa_supplicant successfully backed up")

	f, err := os.OpenFile(networkConfigFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	networkInfo := "network={\n\tssid=\"" + ssid + "\"\n\tpsk=\"" + password + "\"\n}\n"

	if _, err = f.WriteString(networkInfo); err != nil {
		panic(err)
	}
	f.Close()

	resetNetworkCmd := exec.Command("wpa_cli", "reconfigure")
	_, rnErr := resetNetworkCmd.Output()
	if rnErr != nil {
		fmt.Println("Error restarting wpa_supplicant")
		log.Fatal(rnErr)
	}

	// poll status of the wifi connection, early return if success
	for i := 0; i < 20; i++ {
		// check network status
		wifiStatus, wifiStatusErr := getWpaCliStatus()
		if wifiStatusErr != nil {
			log.Fatal(wifiStatusErr)
		}

		if wifiStatus == "COMPLETED" {
			// early return, our work here is done
			return true
		}

		time.Sleep(time.Second)
		fmt.Println("status: " + wifiStatus)
	}

	// reset the wpa_supplicant file to how it was before we tested the new config
	fmt.Println("restoring original contents of wpa_supplicant")
	err = restoreNetworkConfigFile()
	if err != nil {
		log.Fatal("Error restoring network config file.", err)
	}

	fmt.Println("wpa_supplicant has been restored")

	return false
}
