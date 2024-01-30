package subscription

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/LalatinaHub/LatinaSub-go/helper"
)

var (
	blacklistPattern    *regexp.Regexp = regexp.MustCompile("(t.me)")
	subconverterPattern *regexp.Regexp = regexp.MustCompile("(url=|target=)")
	subPattern          *regexp.Regexp = regexp.MustCompile("url=(.+?)(&|$)")

	subList []string = []string{
		"https://raw.githubusercontent.com/LalatinaHub/Mineral/master/result/sub.json",
		"https://raw.githubusercontent.com/mahdibland/V2RayAggregator/master/sub/sub_list.json",
		"https://raw.githubusercontent.com/mfuu/v2ray/master/list.json",
		"https://raw.githubusercontent.com/RenaLio/Mux2sub/main/sub_list",
		"https://raw.githubusercontent.com/RenaLio/Mux2sub/main/urllist",
		"https://beta-test.cloudaccess.host/mylist.json",
	}

	client http.Client = http.Client{
		Timeout: 10 * time.Second,
	}

	subJson SubStruct = SubStruct{
		Id:            0,
		Remarks:       "LalatinaHub/LatinaSub-go",
		Site:          "https://github.com/LalatinaHub/LatinaSub-go",
		Url:           "",
		Update_method: "auto",
		Enabled:       true,
	}

	SubPath string = "subscription/list/"
	SubFile string = SubPath + "sub_list.json"
	subUrls []string
)

func Merge() {
	// Get all the subs url
	for _, subUrl := range subList {
		var (
			buf         *strings.Builder = new(strings.Builder)
			dataTemp    []SubStruct
			subUrlsTemp []string
		)

		resp, err := client.Get(subUrl)
		if err != nil {
			log.Panic(err)
		}
		defer resp.Body.Close()

		if _, err = io.Copy(buf, resp.Body); err != nil {
			log.Panic(err)
		}

		if helper.IsJson(buf.String()) {
			json.Unmarshal([]byte(buf.String()), &dataTemp)
			for _, data := range dataTemp {
				subUrlsTemp = append(subUrlsTemp, strings.Split(data.Url, "|")...)
			}
		} else {
			subUrlsTemp = strings.Split(buf.String(), "\n")
		}

		// Filter
		for _, subUrl := range subUrlsTemp {
			if isMatch := blacklistPattern.MatchString(subUrl); isMatch {
				continue
			} else if isMatch = subconverterPattern.MatchString(subUrl); isMatch {
				subUrl = subPattern.FindStringSubmatch(subUrl)[2]
			}

			if len(subUrl) > 9 {
				if escapedUrl, _ := url.QueryUnescape(subUrl); escapedUrl != "" {
					escapedUrl = strings.Replace(escapedUrl, "http://0.0.0.0:3333/get-base64?content=", "", 10)
					escapedUrl = strings.Replace(escapedUrl, "amp;", "", 10)
					escapedUrl = strings.Replace(escapedUrl, "h312s", "https", 1)
					escapedUrl = strings.Replace(escapedUrl, "httpl.", "http://", 1)
					escapedUrl = strings.TrimRight(escapedUrl, "\r\n")

					if escapedUrl != "" {
						subUrls = append(subUrls, escapedUrl)
					}
				}
			}
		}
	}

	// Filter urls
	filter()

	// Write and populate subs url
	fmt.Println("[+] [Sub] Found", len(subUrls), "subs link !")
	subJson.Url = strings.Join(subUrls, "|")

	out, err := os.Create(SubFile)
	if err != nil {
		log.Panic(err)
	}
	defer out.Close()

	enc := json.NewEncoder(out)
	enc.SetIndent("", "\t")
	enc.SetEscapeHTML(false)

	enc.Encode([]SubStruct{subJson})
}

func Init() {
	// Check and create dir "subscription/list/"
	if _, err := os.Stat(SubPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(SubPath, os.ModePerm)
		} else {
			log.Panic(err)
		}
	}
}

func filter() {
	subUrls = helper.FilterDuplicateString(subUrls)
	subUrls = helper.FilterEmptyString(subUrls)
}
