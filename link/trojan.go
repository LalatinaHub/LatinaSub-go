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

var _ Link = (*Trojan)(nil)

func init() {
	common.Must(RegisterParser(&Parser{
		Name:   "Trojan",
		Scheme: []string{"trojan"},
		Parse: func(u *url.URL) (Link, error) {
			link := &Trojan{}
			return link, link.Parse(u)
		},
	}))
}

// TrojanQt5 represents a parsed Trojan-Qt5 link
type Trojan struct {
	Password        string
	Address         string
	Port            uint16
	Remarks         string
	AllownInsecure  bool
	TFO             bool
	TLS             bool
	WebsocketPath   string
	GrpcServiceName string
	Host            string
	ServerName      string
	Type            string
}

// Options implements Link
func (l *Trojan) Options() *option.Outbound {
	return &option.Outbound{
		Type: C.TypeTrojan,
		Tag:  l.Remarks,
		TrojanOptions: option.TrojanOutboundOptions{
			ServerOptions: option.ServerOptions{
				Server:     l.Address,
				ServerPort: l.Port,
			},
			Password: l.Password,
			TLS: &option.OutboundTLSOptions{
				Enabled:    l.TLS,
				ServerName: l.Address,
				Insecure:   l.AllownInsecure,
			},
			DialerOptions: option.DialerOptions{
				TCPFastOpen: l.TFO,
			},
			Transport: &option.V2RayTransportOptions{
				Type: l.Type,
				WebsocketOptions: option.V2RayWebsocketOptions{
					Path: l.WebsocketPath,
					Headers: map[string]string{
						"Host": l.Host,
					},
				},
				GRPCOptions: option.V2RayGRPCOptions{
					ServiceName: l.GrpcServiceName,
				},
				QUICOptions: option.V2RayQUICOptions{},
			},
		},
	}
}

// Parse implements Link
//
// trojan://password@domain:port?allowinsecure=value&tfo=value#remarks
func (l *Trojan) Parse(u *url.URL) error {
	if u.Scheme != "trojan" {
		return E.New("not a trojan-qt5 link")
	}
	port, err := strconv.ParseUint(u.Port(), 10, 16)
	if err != nil {
		return E.Cause(err, "invalid port")
	}
	l.Address = u.Hostname()
	l.Port = uint16(port)
	l.Remarks = u.Fragment
	if uname := u.User.Username(); uname != "" {
		l.Password = uname
	}

	queries := u.Query()
	for key, values := range queries {
		switch strings.ToLower(key) {
		case "allowinsecure":
			switch values[0] {
			case "0":
				l.AllownInsecure = false
			default:
				l.AllownInsecure = true
			}
		case "tfo":
			switch values[0] {
			case "0":
				l.TFO = false
			default:
				l.TFO = true
			}
		case "security":
			switch values[0] {
			case "":
				l.TLS = false
			default:
				l.TLS = true
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
			l.ServerName = values[0]
		case "path":
			l.WebsocketPath = values[0]
		case "servicename":
			l.GrpcServiceName = values[0]
		}
	}

	return nil
}
