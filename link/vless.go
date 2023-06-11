package link

import (
	"net/url"
	"strconv"
	"strings"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
)

var _ Link = (*Vless)(nil)

func init() {
	common.Must(RegisterParser(&Parser{
		Name:   "VLESS",
		Scheme: []string{"vless"},
		Parse: func(u *url.URL) (Link, error) {
			link := &Vless{}
			return link, link.Parse(u)
		},
	}))
}

// Vless represents a parsed Vless link
type Vless struct {
	UUID            string
	Address         string
	Port            uint16
	Remarks         string
	AllowInsecure   bool
	TLS             bool
	Reality         bool
	TransportPath   string
	GrpcServiceName string
	Host            string
	SNI             string
	Type            string
	PubKey          string
	ShortId         string
}

// Options implements Link
func (l *Vless) Options() *option.Outbound {
	out := &option.Outbound{
		Type: C.TypeVLESS,
		Tag:  l.Remarks,
		VLESSOptions: option.VLESSOutboundOptions{
			ServerOptions: option.ServerOptions{
				Server:     l.Address,
				ServerPort: l.Port,
			},
			UUID: l.UUID,
		},
	}

	if l.TLS {
		out.VLESSOptions.TLS = &option.OutboundTLSOptions{
			Enabled:    l.TLS,
			ServerName: l.SNI,
			Insecure:   l.AllowInsecure,
			DisableSNI: false,
		}

		if l.Reality {
			out.VLESSOptions.TLS.UTLS = &option.OutboundUTLSOptions{
				Enabled: true,
			}
			out.VLESSOptions.TLS.Reality = &option.OutboundRealityOptions{
				Enabled:   true,
				PublicKey: l.PubKey,
				ShortID:   l.ShortId,
			}
		}
	}

	transport := &option.V2RayTransportOptions{
		Type: l.Type,
	}

	switch l.Type {
	case C.V2RayTransportTypeHTTP:
		transport.HTTPOptions.Path = l.TransportPath
		if l.Host != "" {
			transport.HTTPOptions.Host = []string{l.Host}
			if transport.HTTPOptions.Headers == nil {
				transport.HTTPOptions.Headers = map[string]option.Listable[string]{}
			}
			transport.HTTPOptions.Headers["Host"] = option.Listable[string]{l.Host}
		}
	case C.V2RayTransportTypeWebsocket:
		if l.TransportPath == "" {
			l.TransportPath = "/"
		}
		transport.WebsocketOptions.Path = l.TransportPath
		transport.WebsocketOptions.Headers = map[string]option.Listable[string]{
			"Host": {l.Host},
		}
	case C.V2RayTransportTypeQUIC:
		// do nothing
	case C.V2RayTransportTypeGRPC:
		transport.GRPCOptions.ServiceName = l.GrpcServiceName
	default:
		transport = nil
	}

	out.VLESSOptions.Transport = transport

	return out
}

// Parse implements Link
//
// vless://password@domain:port?allowinsecure=value&tfo=value#remarks
func (l *Vless) Parse(u *url.URL) error {

	if u.Scheme != "vless" {
		return E.New("not a vless link")
	}
	port, err := strconv.ParseUint(u.Port(), 10, 16)
	if err != nil {
		return E.Cause(err, "invalid port")
	}
	l.Address = u.Hostname()
	l.Port = uint16(port)
	l.AllowInsecure = true
	l.Remarks = u.Fragment
	if uname := u.User.Username(); uname != "" {
		l.UUID = uname
	}

	queries := u.Query()
	for key, values := range queries {
		switch strings.ToLower(key) {
		case "allowinsecure":
			switch values[0] {
			case "0":
				l.AllowInsecure = false
			}
		case "security":
			l.TLS = false
			switch values[0] {
			case "tls":
				l.TLS = true
			case "reality":
				l.TLS = true
				l.Reality = true
			}
		case "type":
			switch values[0] {
			case "tcp":
			default:
				l.Type = values[0]
			}
		case "host":
			l.Host = values[0]
		case "sni":
			l.SNI = values[0]
		case "path":
			l.TransportPath = values[0]
		case "servicename":
			l.GrpcServiceName = values[0]
		case "pbk":
			l.PubKey = values[0]
		case "sid":
			l.ShortId = values[0]
		}
	}

	return nil
}
