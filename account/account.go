package account

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/LalatinaHub/LatinaSub-go/helper"
	"github.com/LalatinaHub/LatinaSub-go/link"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

var (
	cdnHost string = "speed.cloudflare.com"
	sniHost string = "safebrowsing.google.com"
)

type Account struct {
	Link     string
	Outbound option.Outbound
}

func New(link string) *Account {
	account := Account{Link: link}
	account.Outbound = account.buildOutbound()

	return &account
}

func (account *Account) buildOutbound() option.Outbound {
	defer helper.CatchError(true)

	var outbound option.Outbound
	if parsedNode, _ := url.Parse(account.Link); parsedNode != nil {
		if val, err := link.Parse(parsedNode); val != nil {
			outbound = option.Outbound{
				Type: val.Options().Type,
				Tag:  val.Options().Tag,
			}

			switch val.Options().Type {
			case C.TypeVMess:
				outbound.VMessOptions = val.Options().VMessOptions
			case C.TypeVLESS:
				outbound.VLESSOptions = val.Options().VLESSOptions
			case C.TypeTrojan:
				outbound.TrojanOptions = val.Options().TrojanOptions
			case C.TypeShadowsocks:
				outbound.ShadowsocksOptions = val.Options().ShadowsocksOptions
			case C.TypeShadowsocksR:
				outbound.ShadowsocksROptions = val.Options().ShadowsocksROptions
			}
		} else if err != nil {
			fmt.Println("[Error]", err.Error())
		}
	}

	return outbound
}

func (account Account) PopulateCDN() *option.Outbound {
	switch account.Outbound.Type {
	case C.TypeVMess:
		account.Outbound.VMessOptions.Server = cdnHost
	case C.TypeVLESS:
		account.Outbound.VLESSOptions.Server = cdnHost
	case C.TypeTrojan:
		account.Outbound.TrojanOptions.Server = cdnHost
	case C.TypeShadowsocks:
		account.Outbound.ShadowsocksOptions.Server = cdnHost
	case C.TypeShadowsocksR:
		account.Outbound.ShadowsocksROptions.Server = cdnHost
	}

	return &account.Outbound
}

func (account Account) PopulateSNI() *option.Outbound {
	var TLS *option.OutboundTLSOptions

	switch account.Outbound.Type {
	case C.TypeVMess:
		TLS = account.Outbound.VMessOptions.TLS
	case C.TypeVLESS:
		TLS = account.Outbound.VLESSOptions.TLS
	case C.TypeTrojan:
		TLS = account.Outbound.TrojanOptions.TLS
	case C.TypeShadowsocks, C.TypeShadowsocksR:
		var (
			params        []string
			param, param1 string
			obfs          string = "http"
		)

		// Get the parameter
		switch account.Outbound.Type {
		case C.TypeShadowsocks:
			params = strings.Split(account.Outbound.ShadowsocksOptions.PluginOptions, ";")
			param1 = account.Outbound.ShadowsocksOptions.PluginOptions

		case C.TypeShadowsocksR:
			params = strings.Split(account.Outbound.ShadowsocksROptions.ObfsParam, ";")
			param1 = account.Outbound.ShadowsocksOptions.PluginOptions
		}

		if m, _ := regexp.MatchString("tls", param1); m {
			obfs = "tls"
		}

		// Loop trough plugin-opts and change the host while left other untouched
		for _, parts := range params {
			vals := strings.Split(parts, "=")

			if strings.HasSuffix(vals[0], "host") {
				param = param + vals[0] + "=" + sniHost + ";"
				continue
			}
			param = param + strings.Join(vals, "=") + ";"
		}

		// Remove ; at the end of opts
		param = strings.TrimSuffix(param, ";")

		// Re-assign parameter
		switch account.Outbound.Type {
		case C.TypeShadowsocks:
			account.Outbound.ShadowsocksOptions.PluginOptions = param
		case C.TypeShadowsocksR:
			account.Outbound.ShadowsocksROptions.ObfsParam = fmt.Sprintf("obfs=%s;obfs-host=%s", obfs, sniHost)
		}

		return &account.Outbound
	default:
		return &account.Outbound
	}

	// Make sure TLS is assigned
	if TLS != nil {
		*TLS = option.OutboundTLSOptions{
			Enabled:    TLS.Enabled,
			DisableSNI: false,
			ServerName: sniHost,
			Insecure:   true,
		}
	}

	return &account.Outbound
}
