package device

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
)

type D interface {
	Initialize() error
	StartSlideshow() error
	TestWifiConfiguration(ssid string, password string) bool
}

type Device struct {
	SlideshowInterval int
}

/* Production Device with member functions that interface with a real device */
type ProdDevice struct {
	Device
}

/* Development variant of Device to be used for testing and validation */
type DevDevice struct {
	Device
}

func NewProdDevice(slideshowInterval int) *ProdDevice {
	prodDevice := ProdDevice{}
	prodDevice.SlideshowInterval = slideshowInterval
	return &prodDevice
}

func NewDevDevice(slideshowInterval int) *DevDevice {
	devDevice := DevDevice{}
	devDevice.SlideshowInterval = slideshowInterval
	return &devDevice
}

func getWpaCliStatus() (string, error) {
	// check network status
	checkWifiStatus := exec.Command("wpa_cli", "status")
	wifiStatusStdout, wifiStatusErr := checkWifiStatus.CombinedOutput()
	if wifiStatusErr != nil {
		return "", errors.New("Error checking wpa_cli status")
	}

	// Extract the wpa_state (status)
	regExStatus := regexp.MustCompile(`wpa_state=(.+)\n`)
	// possible states (may not be exhaustive): "ASSOCIATING", "SCANNING", "COMPLETED", "DISCONNECTED", "INACTIVE"
	m := regExStatus.FindStringSubmatch(string(wifiStatusStdout[:]))

	if len(m) == 0 {
		return "", errors.New("wpa_state not found in wpa_cli status")
	}

	// the first match is stored at index 0, the first actual substring is stored at index 1
	wifiStatus := m[1]

	return wifiStatus, nil
}

const networkConfigFilename = "/etc/wpa_supplicant/wpa_supplicant-wlan0.conf"
const networkConfigFilenameBackup = "/etc/wpa_supplicant/wpa_supplicant-wlan0.conf.orig"

func backupNetworkConfigFile() error {
	cmdString := fmt.Sprintf("sudo cp %v %v", networkConfigFilename, networkConfigFilenameBackup)

	return exec.Command("bash", "-c", cmdString).Run()
}

func restoreNetworkConfigFile() error {
	cmdString := fmt.Sprintf("sudo cp %v %v", networkConfigFilenameBackup, networkConfigFilename)

	return exec.Command("bash", "-c", cmdString).Run()
}
