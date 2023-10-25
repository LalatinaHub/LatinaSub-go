package main

import (
	"fmt"
	"os"
	"time"

	latinasub "github.com/LalatinaHub/LatinaSub-go"
)

var (
	saveToDB = false
)

func main() {
	// Start the main func
	for _, arg := range os.Args {
		switch arg {
		case "save_to_db":
			saveToDB = true
		}
	}

	latinasub.Start([]string{}, saveToDB)

	fmt.Printf("\n==============================\n")
	fmt.Println("Software will exit in 10 second !")
	time.Sleep(10 * time.Second)
}
