package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/LalatinaHub/LatinaSub-go/provider"
	"github.com/LalatinaHub/LatinaSub-go/subscription"
	"github.com/sagernet/sing-box/option"
)

var (
	subconverterUrl string   = "https://api.v1.mk/"
	channels        chan int = make(chan int, 20) // Less amount of concurrent seems give better result and stability
	maxNodes        int      = -1
	wg              sync.WaitGroup

	subFile []subscription.SubStruct

	client http.Client = http.Client{
		Timeout: 30 * time.Second,
	}
)

func worker(subUrl string) []option.Outbound {
	var (
		nodes                  = []option.Outbound{}
		buf   *strings.Builder = new(strings.Builder)
	)

	resp, err := client.Get(subUrl)
	if err != nil || resp.StatusCode != 200 {
		return nodes
	}

	// We assume fetch success
	io.Copy(buf, resp.Body)

	content := provider.DecodeBase64Safe(buf.String())
	for _, l := range strings.Split(content, "\n") {
		htmlTagPattern := regexp.MustCompile(`<[\w(\/?)]+>`)
		for _, k := range htmlTagPattern.Split(l, -1) {
			if strings.HasPrefix(k, "http") {
				nodes = append(nodes, worker(k)...)
			}
		}
	}

	outbounds, err := provider.Parse(content)
	if err != nil {
		fmt.Println("[-] Provider:", err.Error())
	}

	for _, outbound := range outbounds {
		if _, err = json.Marshal(outbound); err != nil {
			fmt.Println("[-] Invalid:", err)
			fmt.Println("[-] Parsing:", content)
		} else {
			nodes = append(nodes, outbound)
		}
	}

	if len(nodes) == 0 && strings.HasPrefix(subUrl, subconverterUrl) {
		for _, i := range strings.Split(subUrl, "?") {
			for _, x := range strings.Split(i, "&") {
				if strings.HasPrefix(x, "url=") {
					nodes = append(nodes, worker(strings.Split(x, "=")[1])...)
				}
			}
		}
	}

	return nodes
}

func Run() []option.Outbound {
	var nodes []option.Outbound
	file, _ := os.ReadFile(subscription.SubFile)
	json.Unmarshal(file, &subFile)

	for _, subData := range subFile {
		subUrls := strings.Split(subData.Url, "|")

		for i, subUrl := range subUrls {
			if !strings.HasPrefix(subUrl, "http") {
				continue
			}

			// Limiter
			if len(nodes) > maxNodes && maxNodes > 0 {
				break
			}

			wg.Add(1)

			// Guard
			// Code will be blocked here if there's no empty space
			channels <- i

			go func(id int, subUrl string) {
				defer func() {
					wg.Done()

					// Release id to make an empty space
					<-channels
				}()

				nodes = append(nodes, worker(subUrl)...)
				fmt.Printf("[+] [%d/%d] %s : %d\n", id+1, len(subUrls), subUrl, len(nodes))
			}(i, subUrl)
		}
	}
	wg.Wait()

	return nodes
}
