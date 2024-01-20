package device

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
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
	SlideshowInterval int
	shutdownServer    func() error
}

/* Production Device with member functions that interface with a real device */
type ProdDevice struct {
	Device
}

func NewProdDevice(slideshowInterval int) *ProdDevice {
	prodDevice := ProdDevice{}
	prodDevice.SlideshowInterval = slideshowInterval
	return &prodDevice
}

/* Development variant of Device to be used for testing and validation */
type DevDevice struct {
	Device
}

func NewDevDevice(slideshowInterval int) *DevDevice {
	devDevice := DevDevice{}
	devDevice.SlideshowInterval = slideshowInterval
	return &devDevice
}

// Use the FrameBufferImage(viewer) to run the slideshow on the attached screen
func (d *ProdDevice) StartSlideshow() bool {
	fmt.Println("StartSlideshow called (via ProdDevice interface)")

	wd, wdErr := os.Getwd()
	if wdErr != nil {
		log.Fatal(wdErr)
	}

	// launch the FrameBuffer Image Viewer from a virtual terminal to screen 1, randomize order, auto-zoom, set display interval to 5sec, hide the image metadata
	commandString := "sudo fbi -d /dev/fb0 -T 1 -a -u -t " + fmt.Sprint(d.SlideshowInterval) + " --blend 250 --noverbose ~/photos/*.[^mM4]*"

	fmt.Printf("cwd: %v\n", wd)
	cmd := exec.Command("bash", "-c", commandString)
	fmt.Printf("command: %+v\n", cmd)
	output, err := cmd.Output()
	fmt.Printf(string(output[:]))
	if err != nil {
		log.Fatal(err)
	}
	return true
}

// Use the FrameBufferImage(viewer) to run the slideshow on the attached screen
func (d *DevDevice) StartSlideshow() bool {
	fmt.Println("StartSlideshow called (via DevDevice interface)")
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
	return true
}

// Deactivate the Raspberry Pi's Access Point mode
func (d *Device) DeactivateAccessPointMode() bool {
	fmt.Println("DeactivateAccessPointMode called")
	return true
}

// Start hosting simple web server
func (d *Device) InitializeWebServer(credentialChannel chan map[string]string) bool {
	fmt.Println("InitializeWebServer called (via Device interface)")
	return true
}

// Start hosting simple web server
func (d *Device) InitializeWebServer(credentialChannel chan map[string]string) bool {
	fmt.Println("InitializeWebServer called (via Device interface)")

	sm := http.NewServeMux()

	sm.HandleFunc("/wifi", formSubmissionHandler(credentialChannel))

	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				sm.ServeHTTP(w, r)
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
	fmt.Println("KillWebServer called (via Device interface)")
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

func formSubmissionHandler(credentialChannel chan map[string]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
			fmt.Fprintf(w, "Credentials received:\nssid: %v, password: %v", ssid[0], password[0])

			credentials := make(map[string]string)
			credentials["ssid"] = ssid[0]
			credentials["password"] = password[0]

			// push the credentials we received into the credential channel for consumption elsewhere in the app
			credentialChannel <- credentials
			w.WriteHeader(http.StatusProcessing)
			return
		}

		w.WriteHeader(http.StatusTeapot)
		fmt.Fprintf(w, "Missing SSID or Password")
	}
}
