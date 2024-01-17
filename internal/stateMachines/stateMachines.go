package stateMachines

import (
	"context"
	"fmt"
	"log"
	
	"main/internal/device"

	"github.com/looplab/fsm"
)

// initialize finite state machine to model the application flow
func InitWifiOnboardingWizard(d device.D) *fsm.FSM {
	machine := fsm.NewFSM(
		"init",
		fsm.Events{
			{Name: "next", Src: []string{"init"}, Dst: "activateAccessPointMode"},
			{Name: "next", Src: []string{"activateAccessPointMode"}, Dst: "configureCaptivePortal"},
			{Name: "next", Src: []string{"configureCaptivePortal"}, Dst: "startServer"},
			{Name: "next", Src: []string{"startServer"}, Dst: "displayWifiSetupInstructions"},
			{Name: "next", Src: []string{"displayWifiSetupInstructions"}, Dst: "waitingForCredentials"},
			{Name: "retry", Src: []string{"waitingForCredentials"}, Dst: "waitingForCredentials"},
			{Name: "next", Src: []string{"waitingForCredentials"}, Dst: "validateWifiCredentials"},
			{Name: "failure", Src: []string{"validateWifiCredentials"}, Dst: "waitingForCredentials"},
			{Name: "success", Src: []string{"validateWifiCredentials"}, Dst: "killServer"},
			{Name: "next", Src: []string{"killServer"}, Dst: "disableCaptivePortal"},
			{Name: "next", Src: []string{"disableCaptivePortal"}, Dst: "deactivateAccessPointMode"},
			{Name: "next", Src: []string{"deactivateAccessPointMode"}, Dst: "done"},
		},
		fsm.Callbacks{
			"activateAccessPointMode": func(ctx context.Context, e *fsm.Event) {
				fmt.Println("Activating access point mode...")

				success := d.ActivateAccessPointMode()
				if success {
					e.FSM.Event(ctx, "next")
				} else {
					log.Fatal("error activating access point mode")
				}
			},
			"deactivateAccessPointMode": func(ctx context.Context, e *fsm.Event) {
				fmt.Println("Deactivating access point mode...")

				success := d.DeactivateAccessPointMode()
				if success {
					e.FSM.Event(ctx, "next")
				} else {
					log.Fatal("error deactivating access point mode")
				}
			},
			"configureCaptivePortal": func(ctx context.Context, e *fsm.Event) {
				fmt.Println("Configuring captive portal...")

				success := d.EnableCaptivePortal()
				if success {
					e.FSM.Event(ctx, "next")
				} else {
					log.Fatal("error activating captive portal")
				}
			},
			"disableCaptivePortal": func(ctx context.Context, e *fsm.Event) {
				fmt.Println("Disabling captive portal...")

				success := d.DisableCaptivePortal()
				if success {
					e.FSM.Event(ctx, "next")
				} else {
					log.Fatal("error disabling captive portal")
				}
			},
			"startServer": func(ctx context.Context, e *fsm.Event) {
				fmt.Println("Starting server...")

				credentialChannel := make(chan map[string]string, 1)

				e.FSM.SetMetadata("credentialChannel", credentialChannel)

				success := d.InitializeWebServer(credentialChannel)

				if success {
					e.FSM.Event(ctx, "next")
				} else {
					log.Fatal("error starting server")
				}
			},
			"killServer": func(ctx context.Context, e *fsm.Event) {
				fmt.Println("Killing server...")

				error := d.KillWebServer()
				if error != nil {
					log.Fatal("error killing server")
				} else {
					e.FSM.Event(ctx, "next")
				}
			},
			"displayWifiSetupInstructions": func(ctx context.Context, e *fsm.Event) {
				fmt.Println("Displaying WiFi setup instructions...")

				d.DisplayWifiSetupInstructionImage()

				e.FSM.Event(ctx, "next")
			},
			"waitingForCredentials": func(ctx context.Context, e *fsm.Event) {
				fmt.Println("Waiting for credentials from client...")

				credentialChannel, ok := e.FSM.Metadata("credentialChannel")

				if !ok {
					log.Fatal("Error retrieving credential channel")
				}

				credentials := <-credentialChannel.(chan map[string]string)

				ssid, ssidOk := credentials["ssid"]
				password, passwordOk := credentials["password"]

				if ssidOk && passwordOk {
					fmt.Printf("\nCredentials collected:\nssid: %v\npassword: %v\n\n", ssid, password)

					// persist values to metadata for future consumption
					e.FSM.SetMetadata("ssid", ssid)
					e.FSM.SetMetadata("password", password)

					e.FSM.Event(ctx, "next")
				} else {
					fmt.Println("Either the ssid or password is missing")

					// remove metadata keys to ensure removal of stale data
					e.FSM.DeleteMetadata("ssid")
					e.FSM.DeleteMetadata("password")

					e.FSM.Event(ctx, "retry")
				}
			},
			"validateWifiCredentials": func(ctx context.Context, e *fsm.Event) {
				fmt.Println("Validating credentials...")

				s, _ := e.FSM.Metadata("ssid")
				p, _ := e.FSM.Metadata("password")

				ssid, ssidOk := s.(string)
				password, passwordOk := p.(string)

				if !ssidOk || !passwordOk {
					fmt.Println("Either the ssid or password is missing")
					e.FSM.Event(ctx, "failure")
					return
				}

				success := d.TestWifiConfiguration(ssid, password)

				if success {
					fmt.Println("Successfully connected to WiFi :)")
					e.FSM.Event(ctx, "success")
				} else {
					fmt.Println("Invalid credentials, please try again.")
					e.FSM.Event(ctx, "failure")
				}
			},
			"success": func(ctx context.Context, e *fsm.Event) {
				// close the channel used to relay credentials between the server and the state machine logic
				credentialChannel, ok := e.FSM.Metadata("credentialChannel")
				if !ok {
					log.Fatal("credentialChannel not found")
				}

				close(credentialChannel.(chan map[string]string))
			},
			// create a file to indicate that onboarding has been completed
			"finalizeOnboarding": func(ctx context.Context, e *fsm.Event) { fmt.Println("current state: " + e.FSM.Current()) },
		},
	)

	return machine
}
