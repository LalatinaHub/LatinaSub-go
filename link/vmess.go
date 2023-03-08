package link

import (
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

// Vmess is the base struct of vmess link
type Vmess struct {
	Tag              string
	Server           string
	ServerPort       uint16
	UUID             string
	AlterID          int
	Security         string
	Transport        string
	TransportPath    string
	Host             string
	SNI              string
	TLS              bool
	TLSAllowInsecure bool
	Ver              string
}

// Options implements Link
func (v *Vmess) Options() *option.Outbound {
	out := &option.Outbound{
		Type: C.TypeVMess,
		Tag:  v.Tag,
		VMessOptions: option.VMessOutboundOptions{
			ServerOptions: option.ServerOptions{
				Server:     v.Server,
				ServerPort: v.ServerPort,
			},
			UUID:     v.UUID,
			AlterId:  v.AlterID,
			Security: v.Security,
		},
	}

	if v.TLS {
		out.VMessOptions.TLS = &option.OutboundTLSOptions{
			Enabled:    v.TLS,
			Insecure:   v.TLSAllowInsecure,
			ServerName: v.SNI,
			DisableSNI: false,
		}

		if v.SNI == "" {
			out.VMessOptions.TLS.ServerName = v.Host
		}
	}

	transport := &option.V2RayTransportOptions{
		Type: v.Transport,
	}

	switch v.Transport {
	case C.V2RayTransportTypeHTTP:
		transport.HTTPOptions.Path = v.TransportPath
		if v.Host != "" {
			transport.HTTPOptions.Host = []string{v.Host}
			if transport.HTTPOptions.Headers == nil {
				transport.HTTPOptions.Headers = map[string]string{}
			}
			transport.HTTPOptions.Headers["Host"] = v.Host
		}
	case C.V2RayTransportTypeWebsocket:
		transport.WebsocketOptions.Path = v.TransportPath
		transport.WebsocketOptions.Headers = map[string]string{
			"Host": v.Host,
		}
	case C.V2RayTransportTypeQUIC:
		// do nothing
	case C.V2RayTransportTypeGRPC:
		transport.GRPCOptions.ServiceName = v.TransportPath
	default:
		transport = nil
	}

	out.VMessOptions.Transport = transport

	return out
}
