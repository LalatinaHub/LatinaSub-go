package geoip

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	symbolRegex = regexp.MustCompile("[^a-zA-Z0-9 ]")
	ipinfo      = GeoIpJson{}

	IP_RESOLVER_DOMAIN = "https://myip.shylook.workers.dev"
	IP_RESOLVER_PATH   = "/"
)

func Parse(myIp MyIp) GeoIpJson {
	result := GeoIpJson{
		Ip:          myIp.Ip,
		CountryName: "Unknown",
		CountryCode: "XX",
		Region:      "Unknown",
		Org:         "LalatinaHub",
	}

	for _, country := range CountryList {
		if country.Code == myIp.CC {
			result.CountryName = country.Name
			result.CountryCode = country.Code
			result.Region = country.Region
			result.Org = symbolRegex.ReplaceAllString(myIp.Org, "")
			return result
		}
	}

	return result
}

func GetMyIpInfo() GeoIpJson {
	if ipinfo.Ip != "" {
		return ipinfo
	}

	buf := new(strings.Builder)
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := httpClient.Get(IP_RESOLVER_DOMAIN + IP_RESOLVER_PATH)
	if err != nil {
		return ipinfo
	}
	defer resp.Body.Close()

	io.Copy(buf, resp.Body)
	if resp.StatusCode == 200 {
		myIp := MyIp{}
		if err := json.Unmarshal([]byte(buf.String()), &myIp); err == nil {
			ipinfo = Parse(myIp)
		}
	}

	return ipinfo
}
