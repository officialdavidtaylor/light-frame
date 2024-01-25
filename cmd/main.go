package main
import (
	"log"
	"main/internal/sys"
)

func main() {
	// set default values that can be overridden with values saved in config file
	config := sys.Config{
		SlideshowInterval: 5,
	}

	rErr := config.Read()
	if rErr != nil {
		log.Fatal(rErr)
	}
}
