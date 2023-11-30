package sandbox

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/LalatinaHub/LatinaSub-go/account"
	"github.com/LalatinaHub/LatinaSub-go/helper"
	"github.com/LalatinaHub/LatinaSub-go/ipapi"
	B "github.com/sagernet/sing-box"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

var (
	populateType     = []string{"cdn", "sni"}
	connectivityHost = []string{"http://ipapi.co/json", "http://ipinfo.io/json", "http://google.com"}
)

type SandBox struct {
	Outbound    option.Outbound
	ConnectMode []string
	IpapiObj    ipapi.Ipapi
}

func worker(node option.Outbound, connectMode string) (string, ipapi.Ipapi) {
	var (
		acc        = account.New(node)
		options    option.Options
		listenPort uint
	)

	// Guard
	if acc.Outbound.Type == "" {
		return "", ipapi.Ipapi{}
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
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyClient),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	for _, host := range connectivityHost {
		buf := new(strings.Builder)
		resp, err := httpClient.Get(host)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		io.Copy(buf, resp.Body)
		if resp.StatusCode == 200 {
			if strings.HasSuffix(host, "com") {
				return connectMode, ipapi.Parse("{}")
			}
			return connectMode, ipapi.Parse(buf.String())
		}
	}

	return "", ipapi.Ipapi{}
}

func Test(node option.Outbound) *SandBox {
	var sb SandBox = SandBox{}

	// Constructor
	sb.Outbound = node

	for _, t := range populateType {
		switch sb.Outbound.Type {
		case C.TypeShadowsocksR, C.TypeShadowsocks:
			if t == "cdn" {
				continue
			}
		}

		mode, ipapi := func(node option.Outbound, t string) (string, ipapi.Ipapi) {
			defer helper.CatchError(false)
			return worker(node, t)
		}(sb.Outbound, t)

		if mode != "" {
			sb.ConnectMode = append(sb.ConnectMode, mode)
			sb.IpapiObj = ipapi
		}
	}

	return &sb
}
