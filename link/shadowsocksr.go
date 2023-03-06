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

var _ Link = (*ShadowSocksR)(nil)

func init() {
	common.Must(RegisterParser(&Parser{
		Name:   "ShadowsocksR",
		Scheme: []string{"ssr"},
		Parse: func(u *url.URL) (Link, error) {
			link := &ShadowSocksR{}
			return link, link.Parse(u)
		},
	}))
}

// ShadowSocks represents a parsed shadowsocks link
type ShadowSocksR struct {
	Method        string `json:"method,omitempty"`
	Password      string `json:"password,omitempty"`
	Address       string `json:"address,omitempty"`
	Port          uint16 `json:"port,omitempty"`
	Ps            string `json:"ps,omitempty"`
	Obfs          string `json:"obfs,omitempty"`
	ObfsParam     string `json:"obfs_param,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
	ProtocolParam string `json:"protocol_param,omitempty"`
}

// Options implements Link
func (l *ShadowSocksR) Options() *option.Outbound {
	return &option.Outbound{
		Type: C.TypeShadowsocksR,
		Tag:  l.Ps,
		ShadowsocksROptions: option.ShadowsocksROutboundOptions{
			ServerOptions: option.ServerOptions{
				Server:     l.Address,
				ServerPort: l.Port,
			},
			Method:        l.Method,
			Password:      l.Password,
			Obfs:          l.Obfs,
			ObfsParam:     l.ObfsParam,
			Protocol:      l.Protocol,
			ProtocolParam: l.ProtocolParam,
		},
	}
}

// Parse implements Link
// ssr://server:port:protocol:method:obfs:password_base64/?params_base64
func (l *ShadowSocksR) Parse(u *url.URL) error {
	if u.Scheme != "ssr" {
		return E.New("not a ssr link")
	}

	b64 := u.Host + u.Path
	b, err := base64Decode(b64)
	if err != nil {
		return err
	}
	s := strings.SplitN(string(b), "/?", 2)

	var newLink string = "ssr://"
	for i, value := range strings.Split(s[0], ":") {
		switch i {
		case 0:
			newLink = newLink + value // Server
		case 1:
			newLink = newLink + ":" + value // Server Port
		case 2:
			l.Protocol = value // Protocol
		case 3:
			l.Method = value // Method
		case 4:
			l.Obfs = value
		case 5:
			l.Password = doBase64DecodeOrNothing(value) // Password
		}
	}
	u, err = url.Parse(newLink + "/?" + s[1])
	if err != nil {
		return err
	}

	port, err := strconv.ParseUint(u.Port(), 10, 16)
	if err != nil {
		// return E.Cause(err, "invalid port")
		port = 443
	}
	l.Address = u.Hostname()
	l.Port = uint16(port)
	l.Ps = u.Fragment
	queries := u.Query()
	for key, values := range queries {
		switch key {
		case "protoparam":
			l.ProtocolParam = doBase64DecodeOrNothing(values[0])
		case "obfsparam":
			l.ObfsParam = doBase64DecodeOrNothing(values[0])
		case "remarks":
			l.Ps = doBase64DecodeOrNothing(values[0])
		}
	}

	return nil
}
