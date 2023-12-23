package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/looplab/fsm"
)

func startSlideshow() {
	cmd := exec.Command("fbi", "-d", "/dev/fb0", "-T", "1", "-a", "-u", "-t", "5", "--blend", "250", "--noverbose", "~/photos/*.[^mM4]*\"")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Printf("initializing state machine...\n\n")

	// initialize finite state machine to model the application flow
	fsm := fsm.NewFSM(
		"initial",
		fsm.Events{
			{Name: "startWifiOnboardingFlow", Src: []string{"initial"}, Dst: "activateAccessPointMode"},
			// the next few entries describe the Wifi onboarding flow
			{Name: "next", Src: []string{"activateAccessPointMode"}, Dst: "startCaptivePortalService"},
			{Name: "next", Src: []string{"startCaptivePortalService"}, Dst: "displayWifiSetupInstructions"},
			{Name: "testWifiCredentials", Src: []string{"displayWifiSetupInstructions"}, Dst: "validateWifiCredentials"},
			{Name: "failure", Src: []string{"validateWifiCredentials"}, Dst: "displayWifiSetupInstructions"},
			{Name: "success", Src: []string{"validateWifiCredentials"}, Dst: "finalizeOnboarding"},
			{Name: "next", Src: []string{"finalizeOnboarding"}, Dst: "displaySlideshow"},
			// this allows the app to skip the Wifi onboarding flow and cut directly to the slideshow if appropriate
			{Name: "startSlideshow", Src: []string{"initial"}, Dst: "displaySlideshow"},
		},
		fsm.Callbacks{
			"displaySlideshow": func(_ context.Context, e *fsm.Event) {
				startSlideshow()
			},
			"initial":                      func(_ context.Context, e *fsm.Event) { fmt.Println("current state: " + e.FSM.Current()) },
			"activateAccessPointMode":      func(_ context.Context, e *fsm.Event) { fmt.Println("current state: " + e.FSM.Current()) },
			"startCaptivePortalService":    func(_ context.Context, e *fsm.Event) { fmt.Println("current state: " + e.FSM.Current()) },
			"displayWifiSetupInstructions": func(_ context.Context, e *fsm.Event) { fmt.Println("current state: " + e.FSM.Current()) },
			"validateWifiCredentials":      func(_ context.Context, e *fsm.Event) { fmt.Println("current state: " + e.FSM.Current()) },
			// create a file to indicate that onboarding has been completed
			"finalizeOnboarding": func(_ context.Context, e *fsm.Event) { fmt.Println("current state: " + e.FSM.Current()) },
		},
	)

	// check to see if the user has completed onboarding
	if _, err := os.Stat("./build/config/onboardingCompleted"); err == nil {
		// run this block if onboarding has already been completed
		fmt.Println("onboarding has already been completed!")

		// proceed straight to the `startSlideshow` state
		err := fsm.Event(context.Background(), "startSlideshow")
		if err != nil {
			fmt.Println(err)
		}

	} else if errors.Is(err, os.ErrNotExist) {
		// run this block if onboarding has *not* been completed
		fmt.Println("onboarding has not been completed!")

		err := fsm.Event(context.Background(), "startWifiOnboardingFlow")
		if err != nil {
			fmt.Println(err)
		}

		err = fsm.Event(context.Background(), "next")
		if err != nil {
			fmt.Println(err)
		}

		err = fsm.Event(context.Background(), "next")
		if err != nil {
			fmt.Println(err)
		}

		err = fsm.Event(context.Background(), "testWifiCredentials")
		if err != nil {
			fmt.Println(err)
		}

		err = fsm.Event(context.Background(), "failure")
		if err != nil {
			fmt.Println(err)
		}

		err = fsm.Event(context.Background(), "next")
		if err != nil {
			fmt.Println(err)
		}

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		log.Fatal(err)
	}
}
