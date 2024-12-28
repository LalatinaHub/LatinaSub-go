package sandbox

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/LalatinaHub/LatinaSub-go/account"
	"github.com/LalatinaHub/LatinaSub-go/geoip"
	"github.com/LalatinaHub/LatinaSub-go/helper"
	B "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/option"
)

var (
	populateType     = []string{"cdn", "sni"}
	connectivityHost = []string{geoip.IP_RESOLVER_DOMAIN + geoip.IP_RESOLVER_PATH}
)

type SandBox struct {
	Outbound    option.Outbound
	ConnectMode []string
	Geoip       geoip.GeoIpJson
}

func worker(node option.Outbound, connectMode string) (string, geoip.GeoIpJson) {
	var (
		acc        = account.New(node)
		options    option.Options
		listenPort uint
	)

	// Guard
	if acc.Outbound.Type == "" {
		return "", geoip.GeoIpJson{}
	}

	if connectMode == "cdn" {
		options, listenPort = generateConfig(acc.PopulateCDN())
	} else {
		options, listenPort = generateConfig(acc.PopulateSNI())
	}

	box, err := B.New(B.Options{
		Context: context.Background(),
		Options: options,
	})
	if err != nil {
		panic(err)
	}
	defer box.Close()

	// Start sing-box client
	box.Start()

	proxyClient, _ := url.Parse(fmt.Sprintf("socks5://0.0.0.0:%d", listenPort))
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyClient),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	req, _ := http.NewRequest("GET", "https://speed.cloudflare.com/cdn-cgi/trace", nil)
	resp, err := httpClient.Do(req)
	if resp.StatusCode == 200 && err == nil {
		for _, host := range connectivityHost {
			buf := new(strings.Builder)
			req, err = http.NewRequest("GET", host, nil)
			if err != nil {
				panic(err)
			}

			resp, err := httpClient.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			io.Copy(buf, resp.Body)
			if resp.StatusCode == 200 {
				myIp := geoip.MyIp{}
				if err := json.Unmarshal([]byte(buf.String()), &myIp); err == nil {
					return connectMode, geoip.Parse(myIp)
				}
			}
		}
	}

	return "", geoip.GeoIpJson{}
}

func Test(node option.Outbound) *SandBox {
	var sb SandBox = SandBox{}

	// Constructor
	sb.Outbound = node

	for _, t := range populateType {
		mode, geoproxyip := func(node option.Outbound, t string) (string, geoip.GeoIpJson) {
			defer helper.CatchError(false)
			return worker(node, t)
		}(sb.Outbound, t)

		if mode != "" {
			if geoproxyip.Ip != geoip.GetMyIpInfo().Ip {
				sb.ConnectMode = append(sb.ConnectMode, mode)
				sb.Geoip = geoproxyip
			}
		}
	}

	return &sb
}
