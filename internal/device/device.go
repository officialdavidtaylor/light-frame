package device

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"os/exec"
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
func (d *ProdDevice) StartSlideshow() bool {
	// launch the FrameBuffer Image Viewer from a virtual terminal to screen 1, randomize order, auto-zoom, set display interval to 5sec, hide the image metadata
	cmd := exec.Command("fbi", "-d", "/dev/fb0", "-T", "1", "-a", "-u", "-t", fmt.Sprint(d.slideshowTimer), "--blend", "250", "--noverbose", "~/photos/*.[^mM4]*\"")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return true
}

// Use the FrameBufferImage(viewer) to run the slideshow on the attached screen
func (d *DevDevice) StartSlideshow() bool {
	fmt.Println("StartSlideshow called")
	return true
}

// Display Wifi setup instructions with FrameBufferImage(viewer)
func (d *DevDevice) DisplayWifiSetupInstructionImage() bool {
	fmt.Println("DisplayWifiSetupInstructionImage called")
	return true
}

// Display Wifi setup instructions with FrameBufferImage(viewer)
func (d *ProdDevice) DisplayWifiSetupInstructionImage() bool {
	// launch the FrameBuffer Image Viewer from a virtual terminal to screen 1, hide the image metadata
	cmd := exec.Command("fbi", "-d", "/dev/fb0", "-T", "1", "-a", "--noverbose", "./assets/WifiInstructions.png")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
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
	fmt.Println("InitializeWebServer called (via Device interface)")
	return true
}

// Start hosting simple web server
func (d *DevDevice) InitializeWebServer(credentialChannel chan map[string]string) bool {
	fmt.Println("InitializeWebServer called (via DevDevice interface)")

	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				queryParams := r.URL.Query()

				ssid, ssidOk := queryParams["ssid"]
				if !ssidOk {
					fmt.Println("missing ssid query param")
				}

				password, passwordOk := queryParams["password"]
				if !passwordOk {
					fmt.Println("missing password query param")
				}

				if ssidOk && passwordOk {
					fmt.Fprintf(w, "Hello, %q\nssid: %v, password: %v", html.EscapeString(r.URL.Path), ssid[0], password[0])

					credentials := make(map[string]string)
					credentials["ssid"] = ssid[0]
					credentials["password"] = password[0]

					credentialChannel <- credentials
				}

			},
		),
	}

	fmt.Println("starting server")
	go server.ListenAndServe()

	// create server shutdown closure that can be called later
	d.shutdownServer = func() error {
		return server.Shutdown(context.Background())
	}

	return true
}

// Kill simple web server
func (d *Device) KillWebServer() error {
	fmt.Println("KillWebServer called")
	return nil
}

// Kill simple web server
func (d *DevDevice) KillWebServer() error {
	fmt.Println("KillWebServer called")
	return d.shutdownServer()
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
	fmt.Println("TestWifiConfiguration called (via Device interface)")
	return true
}

// Test WiFi configuration
func (d *DevDevice) TestWifiConfiguration(ssid string, password string) bool {
	fmt.Println("TestWifiConfiguration called")

	if ssid == "testSsid" && password == "testPassword" {
		return true
	}

	return false
}
