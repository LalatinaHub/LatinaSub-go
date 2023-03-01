package scraper

import (
	"encoding/base64"
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
	"github.com/LalatinaHub/LatinaSub-go/subscription"
)

var (
	subconverterUrl string   = "https://api.v1.mk/"
	channels        chan int = make(chan int, 20) // Less amount of concurrent seems give better result and stability
	maxNodes        int      = -1
	wg              sync.WaitGroup

	acceptedProtocol []string       = []string{"vmess", "ss", "trojan", "vless" /*,"ssr"*/}
	protocolPattern  *regexp.Regexp = regexp.MustCompile(fmt.Sprintf("(%s)://.+", strings.Join(acceptedProtocol, "|")))
	subFile          []subscription.SubStruct

	client http.Client = http.Client{
		Timeout: 10 * time.Second,
	}
)

func worker(subUrl string) []string {
	var (
		nodes []string
		buf   *strings.Builder = new(strings.Builder)
	)

	resp, err := client.Get(subUrl)
	if err != nil {
		return []string{}
	} else if resp.StatusCode != 200 {
		return []string{}
	}

	// We assume fetch success
	io.Copy(buf, resp.Body)

	// Detech content-type of the successfull fetch
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application") && !strings.HasPrefix(subUrl, subconverterUrl) {
		// Probably type: yaml, x-yaml, octet-stream

		// Re-fetch suburl using subconverter api
		for _, protocolType := range acceptedProtocol {
			switch protocolType {
			case "vmess":
				subUrl = fmt.Sprintf("%ssub?target=v2ray&url=%s", subconverterUrl, subUrl)
			case "vless":
				continue
			default:
				subUrl = fmt.Sprintf("%ssub?target=%s&url=%s", subconverterUrl, protocolType, subUrl)
			}
			nodes = append(nodes, worker(subUrl)...)
		}
	} else {
		str := buf.String()

		// Find all nodes by pattern
		foundNodes := protocolPattern.FindAllString(str, -1)

		// The result maybe base64 encoded, so we try to decode it
		if len(foundNodes) == 0 {
			if decodedStr, _ := base64.StdEncoding.DecodeString(str); decodedStr != nil {
				foundNodes = protocolPattern.FindAllString(string(decodedStr), -1)
			}
		}

		// Populate slices
		nodes = append(nodes, foundNodes...)
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

					// Release id to make a empty space
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
