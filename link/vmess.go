package link

import (
	"encoding/json"
	"net/url"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
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
	TLS              bool
	TLSAllowInsecure bool
	Ver              string
}

type _vmess struct {
	V              number `json:"v,omitempty"`
	Ps             string `json:"ps,omitempty"`
	Add            string `json:"add,omitempty"`
	Port           number `json:"port,omitempty"`
	ID             string `json:"id,omitempty"`
	Aid            number `json:"aid,omitempty"`
	Scy            string `json:"scy,omitempty"`
	Security       string `json:"security,omitempty"`
	SkipCertVerify bool   `json:"skip-cert-verify,omitempty"`
	Net            string `json:"net,omitempty"`
	Type           string `json:"type,omitempty"`
	Host           string `json:"host,omitempty"`
	Path           string `json:"path,omitempty"`
	TLS            string `json:"tls,omitempty"`
	SNI            string `json:"sni,omitempty"`
	ALPN           string `json:"alpn,omitempty"`
}

func init() {
	common.Must(RegisterParser(&Parser{
		Name:   "Vmess",
		Scheme: []string{"vmess"},
		Parse: func(u *url.URL) (Link, error) {
			link := &Vmess{}
			return link, link.Parse(u)
		},
	}))
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
			Enabled:    true,
			Insecure:   v.TLSAllowInsecure,
			ServerName: v.Host,
		}
	}

	opt := &option.V2RayTransportOptions{
		Type: v.Transport,
	}

	switch v.Transport {
	case "":
		opt = nil
	case C.V2RayTransportTypeHTTP:
		opt.HTTPOptions.Path = v.TransportPath
		if v.Host != "" {
			opt.HTTPOptions.Host = []string{v.Host}
			if opt.HTTPOptions.Headers == nil {
				opt.HTTPOptions.Headers = map[string]string{}
			}
			opt.HTTPOptions.Headers["Host"] = v.Host
		}
	case C.V2RayTransportTypeWebsocket:
		opt.WebsocketOptions.Path = v.TransportPath
		opt.WebsocketOptions.Headers = map[string]string{
			"Host": v.Host,
		}
	case C.V2RayTransportTypeQUIC:
		// do nothing
	case C.V2RayTransportTypeGRPC:
		opt.GRPCOptions.ServiceName = v.Host
	}

	out.VMessOptions.Transport = opt
	return out
}

// Parse implements Link
func (l *Vmess) Parse(u *url.URL) error {
	if u.Scheme != "vmess" {
		return E.New("not a vmess link")
	}

	b64 := u.Host + u.Path
	b, err := base64Decode(b64)
	if err != nil {
		return err
	}

	v := _vmess{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	v.Aid = 0
	l.Tag = v.Ps
	l.Server = v.Add
	l.ServerPort = uint16(v.Port)
	l.UUID = v.ID
	l.AlterID = int(v.Aid)
	l.Security = "auto"
	if v.Scy != "" {
		l.Security = v.Scy
	}
	l.Host = v.Host
	l.TransportPath = v.Path
	l.TLS = v.TLS == "tls"
	l.TLSAllowInsecure = v.SkipCertVerify
	// _ = v.Type
	// _ = v.SNI
	// _ = v.ALPN

	switch v.Net {
	case "ws", "websocket":
		l.Transport = C.V2RayTransportTypeWebsocket
	case "http":
		l.Transport = C.V2RayTransportTypeHTTP
	case "grpc":
		l.Transport = C.V2RayTransportTypeGRPC
	}

	return nil
}
