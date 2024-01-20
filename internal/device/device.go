package device

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"time"
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
	commandString := "sudo fbi -d /dev/fb0 -T 1 -a --noverbose ~/assets/WiFi\\ Setup\\ Instructions.png"
	// launch the FrameBuffer Image Viewer from a virtual terminal to screen 1, hide the image metadata
	cmd := exec.Command("bash", "-c", commandString)
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

	sm.Handle("/", http.FileServer(http.Dir("./frontend/static")))

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
func (d *ProdDevice) TestWifiConfiguration(ssid string, password string) bool {
	fmt.Println("TestWifiConfiguration called (via ProdDevice interface)")
	addNetworkInfoCmdString := "sudo raspi-config nonint do_wifi_ssid_passphrase \"" + ssid + "\" \"" + password + "\""

	addNetworkInfoCmd := exec.Command("bash", "-c", addNetworkInfoCmdString)
	_, addNetworkErr := addNetworkInfoCmd.Output()
	if addNetworkErr != nil {
		fmt.Println("Error adding network info to wpa_supplicant")
		log.Fatal(addNetworkErr)
	}

	restartWpaSupplicantCmdString := "wpa_cli reconfigure"

	restartWpaSupplicantCmd := exec.Command("bash", "-c", restartWpaSupplicantCmdString)
	_, restartWpaSupplicant := restartWpaSupplicantCmd.Output()
	if restartWpaSupplicant != nil {
		fmt.Println("Error restarting wpa_supplicant")
		log.Fatal(restartWpaSupplicant)
	}

	// poll status of the wifi connection for 10 seconds
	for i := 0; i < 10; i++ {
		// check network status
		wifiStatus, wifiStatusErr := getWpaCliStatus()

		if wifiStatusErr != nil {
			log.Fatal(wifiStatusErr)
		}

		if wifiStatus == "COMPLETED" {
			return true
		}

		time.Sleep(time.Second)
		fmt.Println("status: " + wifiStatus)
	}

	return false
}

// Test WiFi configuration
func (d *DevDevice) TestWifiConfiguration(ssid string, password string) bool {
	fmt.Println("TestWifiConfiguration called")

	if ssid == "testSsid" && password == "testPassword" {
		return true
	}

	return false
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
	// possible states (may not be exhaustive): "ASSOCIATING", "SCANNING", "COMPLETED", "DISCONNECTED"
	m := regExStatus.FindStringSubmatch(string(wifiStatusStdout[:]))

	if len(m) == 0 {
		return "", errors.New("wpa_state not found in wpa_cli status")
	}

	// the first match is stored at index 0, the first actual substring is stored at index 1
	wifiStatus := m[1]

	return wifiStatus, nil
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
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("Missing SSID or Password"))
	}
}
