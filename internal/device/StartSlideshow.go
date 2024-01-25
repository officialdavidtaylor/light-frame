package device

import (
	"fmt"
	"os/exec"
)

// fallback for testing on a dev device
func (d *Device) StartSlideshow() error {
	fmt.Println("StartSlideshow called (via Device interface)")
	return nil
}

// Use the FrameBufferImage(viewer) to run the slideshow on the attached screen
func (d *ProdDevice) StartSlideshow() error {
	// launch the FrameBuffer Image Viewer from a virtual terminal to screen 1, randomize order, auto-zoom, set display interval, hide the image metadata
	commandString := "sudo fbi -d /dev/fb0 -T 1 -a -u -t " + fmt.Sprint(d.SlideshowInterval) + " --blend 250 --noverbose ~/photos/*.[^mM4]*"

	cmd := exec.Command("bash", "-c", commandString)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
