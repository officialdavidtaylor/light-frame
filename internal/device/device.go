package device

import (
	"fmt"
)

type D interface {
	StartSlideshow() bool
	DisplayWifiSetupInstructionImage() bool
	TestWifiConfiguration(ssid string, password string) bool
	ActivateAccessPointMode() bool
	DeactivateAccessPointMode() bool
	InitializeWebServer(credentialChannel chan map[string]string) bool
	KillWebServer() error
	EnableCaptivePortal() bool
	DisableCaptivePortal() bool
}

type Device struct {
}

/* Production Device with member functions that interface with a real device */
type ProdDevice struct {
	Device
}

/* Development variant of Device to be used for testing and validation */
type DevDevice struct {
	Device
}

// Use the FrameBufferImage(viewer) to run the slideshow on the attached screen
func (d *Device) StartSlideshow() bool {
	fmt.Println("StartSlideshow called")
	return true
}

// Display Wifi setup instructions with FrameBufferImage(viewer)
func (d *Device) DisplayWifiSetupInstructionImage() bool {
	fmt.Println("DisplayWifiSetupInstructionImage called")
	return true
}

// Activate the Raspberry Pi's Access Point mode
func (d *Device) ActivateAccessPointMode() bool {
	fmt.Println("ActivateAccessPointMode called")
	return false
}

// Deactivate the Raspberry Pi's Access Point mode
func (d *Device) DeactivateAccessPointMode() bool {
	fmt.Println("DeactivateAccessPointMode called")
	return false
}

// Start hosting simple web server
func (d *Device) InitializeWebServer(credentialChannel chan map[string]string) bool {
	fmt.Println("InitializeWebServer called")
	return true
}

// Kill simple web server
func (d *Device) KillWebServer() error {
	fmt.Println("KillWebServer called")
	return nil
}

// Enable captive portal to route all AP requests to our local server
func (d *Device) EnableCaptivePortal() bool {
	fmt.Println("EnableCaptivePortal called")
	return true
}

// Disable captive portal
func (d *Device) DisableCaptivePortal() bool {
	fmt.Println("DisableCaptivePortal called")
	return true
}

// Test WiFi configuration
func (d *Device) TestWifiConfiguration(ssid string, password string) bool {
	fmt.Println("TestWifiConfiguration called")
	return false
}
