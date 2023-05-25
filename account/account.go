package account

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"

	"github.com/LalatinaHub/LatinaSub-go/helper"
	"github.com/LalatinaHub/LatinaSub-go/link"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

var (
	cdnHost string = "172.67.73.39"
	sniHost string = "teams.microsoft.com"
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
	case C.TypeShadowsocks:
		var (
			obfs = "tls"
			port = int64(account.Outbound.ShadowsocksOptions.ServerPort)
		)

		if m, _ := regexp.MatchString("80|88", strconv.FormatInt(port, 10)); m {
			obfs = "http"
		}

		account.Outbound.ShadowsocksOptions.Plugin = "obfs-local"
		account.Outbound.ShadowsocksOptions.PluginOptions = fmt.Sprintf("obfs=%s;obfs-host=%s", obfs, sniHost)
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
