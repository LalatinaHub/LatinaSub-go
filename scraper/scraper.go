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

	"github.com/LalatinaHub/LatinaSub-go/helper"
	"github.com/LalatinaHub/LatinaSub-go/link"
	"github.com/LalatinaHub/LatinaSub-go/subscription"
)

var (
	subconverterUrl string   = "https://api.v1.mk/"
	channels        chan int = make(chan int, 20) // Less amount of concurrent seems give better result and stability
	maxNodes        int      = -1
	wg              sync.WaitGroup

	acceptedProtocol []string       = []string{"vmess", "ss", "trojan", "vless"}
	protocolPattern  *regexp.Regexp = regexp.MustCompile(fmt.Sprintf("(%s)://.+", strings.Join(acceptedProtocol, "|")))
	subFile          []subscription.SubStruct

	client http.Client = http.Client{
		Timeout: 30 * time.Second,
	}
)

func worker(subUrl string) []string {
	var (
		nodes                  = []string{}
		buf   *strings.Builder = new(strings.Builder)
	)

	resp, err := client.Get(subUrl)
	if err != nil || resp.StatusCode != 200 {
		return nodes
	}

	// We assume fetch success
	io.Copy(buf, resp.Body)

	str := link.DoBase64DecodeOrNothing(buf.String())
	foundNodes := protocolPattern.FindAllString(str, -1)

	if len(foundNodes) > 0 {
		nodes = append(nodes, foundNodes...)
	} else if !strings.HasSuffix(subUrl, "fool=1") {
		url := subUrl
		if strings.HasPrefix(subUrl, subconverterUrl) {
			for _, i := range strings.Split(subUrl, "?") {
				for _, x := range strings.Split(i, "&") {
					if strings.HasPrefix(x, "url=") {
						url = strings.Split(x, "=")[1]
					}
				}
			}
		}

		for _, protocolType := range acceptedProtocol {
			switch protocolType {
			case "vmess":
				subUrl = fmt.Sprintf("%ssub?target=v2ray&url=%s&fool=1", subconverterUrl, url)
			case "vless":
				continue
			default:
				subUrl = fmt.Sprintf("%ssub?target=%s&url=%s&fool=1", subconverterUrl, protocolType, url)
			}
			nodes = append(nodes, worker(subUrl)...)
		}
	}
	return nodes
}

func Run() []string {
	var nodes []string
	file, _ := os.ReadFile(subscription.SubFile)
	json.Unmarshal(file, &subFile)

	for _, subData := range subFile {
		subUrls := strings.Split(subData.Url, "|")

		for i, subUrl := range subUrls {
			// Limiter
			if len(nodes) > maxNodes && maxNodes > 1 {
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
				fmt.Printf("[%d/%d] %s : %d\n", id+1, len(subUrls), subUrl, len(nodes))
			}(i, subUrl)
		}
	}
	wg.Wait()

	// Filter nodes
	fmt.Println("[filter] Filtering nodes !")
	nodes = filter(nodes)

	return nodes
}

func filter(nodes []string) []string {
	nodes = helper.FilterDuplicateString(nodes)
	nodes = helper.FilterEmptyString(nodes)

	return nodes
}
