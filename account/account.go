package account

import (
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

var (
	cdnHost string = "172.67.73.39"
	sniHost string = "meet.google.com"
)

type Account struct {
	Outbound option.Outbound
}

func New(node option.Outbound) *Account {
	account := Account{Outbound: node}

	return &account
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
