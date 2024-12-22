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
		DNS: &option.DNSOptions{
			Servers: []option.DNSServerOptions{
				{
					Address: "1.1.1.1",
					Detour:  "direct",
				},
			},
		},
		Inbounds: []option.Inbound{
			{
				Type: C.TypeMixed,
				MixedOptions: option.HTTPMixedInboundOptions{
					ListenOptions: option.ListenOptions{
						Listen:     option.NewListenAddress(netip.IPv4Unspecified()),
						ListenPort: uint16(listenPort),
					},
				},
			},
		},
		Outbounds: []option.Outbound{
			*out,
			{
				Tag:  C.TypeDirect,
				Type: C.TypeDirect,
			},
		},
		Route: &option.RouteOptions{
			Rules: []option.Rule{
				{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultRule{
						Protocol: option.Listable[string]{"dns"},
						Outbound: "direct",
					},
				},
				{
					Type: C.RuleTypeDefault,
					DefaultOptions: option.DefaultRule{
						Network:  option.Listable[string]{"udp"},
						Outbound: "direct",
					},
				},
			},
			Final: out.Tag,
		},
	}

	return options, listenPort
}
