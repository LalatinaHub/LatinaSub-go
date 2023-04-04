package latinasub

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/LalatinaHub/LatinaSub-go/blacklist"
	D "github.com/LalatinaHub/LatinaSub-go/db"
	"github.com/LalatinaHub/LatinaSub-go/helper"
	"github.com/LalatinaHub/LatinaSub-go/sandbox"
	"github.com/LalatinaHub/LatinaSub-go/scraper"

	"github.com/LalatinaHub/LatinaSub-go/subscription"
)

var (
	loc, _         = time.LoadLocation("Asia/Jakarta")
	Concurrent int = runtime.NumCPU() * 100
	wg         sync.WaitGroup
	GoodBoxes  []*sandbox.SandBox
)

func initAll() {
	subscription.Init()
	blacklist.Init()
}

func Start(nodes []string) int {
	start := time.Now()

	if concurrentStr, isSet := os.LookupEnv("CONCURRENT"); isSet {
		if cr, _ := strconv.Atoi(concurrentStr); cr > 0 {
			Concurrent = cr
		}
	}

	fmt.Println("Total concurrent:", Concurrent)

	// Initialize all required modules
	initAll()
	ch := make(chan int, Concurrent)
	db := D.New()

	// Scrape nodes from sub list if parameter is empty
	if len(nodes) <= 0 {
		// Merge sub list
		subscription.Merge()

		nodes = scraper.Run()
	}
	numNodes := len(nodes)

	for i, node := range nodes {
		fmt.Println("Testing node no", i, "/", len(nodes))

		// Blacklist guard
		if blacklist.Find(node) {
			continue
		}

		wg.Add(1)

		ch <- 1
		go func(node string, numNodes, id int) {
			// Catch on error and done wg
			defer helper.CatchError(false)

			// Make sure wg is done and channel released
			defer func() {
				wg.Done()
				<-ch
			}()

			// Start the test
			if box := sandbox.Test(node); box != nil {
				if len(box.ConnectMode) > 0 {
					GoodBoxes = append(GoodBoxes, box)
					fmt.Printf("[%d/%d] Connected in %s mode => %d\n", id, numNodes, strings.Join(box.ConnectMode, " and "), len(GoodBoxes))
				} else {
					// Blacklist unused node
					// if !blacklist.Find(node) {
					// 	blacklist.Save(node)
					// 	fmt.Printf("[%d/%d] Blacklisted\n", id, numNodes)
					// }
				}
			}
		}(node, numNodes, i)
	}

	// Wait for all process
	wg.Wait()

	// Write blacklist
	fmt.Println("Saving blacklisted nodes, please wait !")
	blacklist.Write()

	// Write all result to database
	fmt.Println("Saving result to database, please wait !")
	db.FlushAndCreate()
	db.Save(GoodBoxes)

	// Log collapsed time
	fmt.Println("Total CPU:", runtime.NumCPU())
	fmt.Println("Total time collapsed:", time.Since(start))
	fmt.Println("Total accounts:", db.TotalAccount)
	fmt.Println("Finish Time:", time.Now().In(loc).Format("2006-01-02 15:04:05"))

	return db.TotalAccount
}
