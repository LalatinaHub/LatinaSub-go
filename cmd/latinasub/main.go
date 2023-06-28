package main

import (
	"fmt"
	"time"

	latinasub "github.com/LalatinaHub/LatinaSub-go"
)

func main() {
	// Start the main func
	latinasub.Start([]string{}, false)

	fmt.Printf("\n==============================\n")
	fmt.Println("Software will exit in 10 second !")
	time.Sleep(10 * time.Second)
}
