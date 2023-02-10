package main

import (
	"fmt"
	"time"

	latinasub "github.com/LalatinaHub/LatinaSub-go"
)

func main() {
	// Info
	fmt.Println("Total Concurrent:", latinasub.Concurrent)

	// Start the main func
	totalAccount := latinasub.Start()

	fmt.Printf("\n==============================\n")
	fmt.Println("Result:", totalAccount)

	fmt.Println("Software will exit in 10 second !")
	time.Sleep(10 * time.Second)
}
