package link

import (
	"encoding/json"
	"net/url"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
)

func init() {
	common.Must(RegisterParser(&Parser{
		Name:   "V2RayNG",
		Scheme: []string{"vmess"},
		Parse: func(u *url.URL) (Link, error) {
			link := &VMessV2RayNG{}
			return link, link.Parse(u)
		},
	}))
}

// VMessV2RayNG is the vmess link of V2RayNG
type VMessV2RayNG struct {
	Vmess

	Ver string
}

type _vmessV2RayNG struct {
	V              number `json:"v,omitempty"`
	Ps             string `json:"ps,omitempty"`
	Add            string `json:"add,omitempty"`
	Port           number `json:"port,omitempty"`
	ID             string `json:"id,omitempty"`
	Aid            number `json:"aid,omitempty"`
	Scy            string `json:"scy,omitempty"`
	Security       string `json:"security,omitempty"`
	SkipCertVerify any    `json:"skip-cert-verify,omitempty"` // Have 2 possibly data type
	Net            string `json:"net,omitempty"`
	Type           string `json:"type,omitempty"`
	Host           string `json:"host,omitempty"`
	Path           string `json:"path,omitempty"`
	TLS            string `json:"tls,omitempty"`
	SNI            string `json:"sni,omitempty"`
	ALPN           string `json:"alpn,omitempty"`
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

	v := _vmessV2RayNG{}
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
	if v.SkipCertVerify != nil {
		l.TLSAllowInsecure = true
	} else {
		l.TLSAllowInsecure = false
	}
	l.SNI = v.SNI
	// _ = v.Type
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
