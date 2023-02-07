package sandbox

import (
	"net/netip"

	"github.com/LalatinaHub/LatinaSub-go/helper"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func generateConfig(out *option.Outbound) (option.Options, uint) {
	listenPort := helper.GetFreePort()
	options := option.Options{
		Log: &option.LogOptions{
			Disabled:  true,
			Level:     "error",
			Timestamp: true,
		},
		Inbounds: []option.Inbound{
			{
				Type: C.TypeHTTP,
				HTTPOptions: option.HTTPMixedInboundOptions{
					ListenOptions: option.ListenOptions{
						Listen:     option.ListenAddress(netip.IPv4Unspecified()),
						ListenPort: uint16(listenPort),
					},
				},
			},
		},
		Outbounds: []option.Outbound{
			*out,
		},
	}

	return options, listenPort
}
