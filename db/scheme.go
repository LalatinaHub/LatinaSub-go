package db

// Merge from SS, SSR, Vmess, Vless/Trojan
type DBScheme struct {
	Server        string `json:"server,omitempty"`         // 0
	Ip            string `json:"ip,omitempty"`             // 1
	ServerPort    int    `json:"server_port,omitempty"`    // 2
	UUID          string `json:"uuid,omitempty"`           // 3
	Password      string `json:"password,omitempty"`       // 4
	Security      string `json:"security,omitempty"`       // 5
	AlterId       int    `json:"alter_id,omitempty"`       // 6
	Method        string `json:"method,omitempty"`         // 7
	Plugin        string `json:"plugin,omitempty"`         // 8
	PluginOpts    string `json:"plugin_opts,omitempty"`    // 9
	Protocol      string `json:"protocol,omitempty"`       // 10
	ProtocolParam string `json:"protocol_param,omitempty"` // 11
	OBFS          string `json:"obfs,omitempty"`           // 12
	OBFSParam     string `json:"obfs_param,omitempty"`     // 13
	Host          string `json:"host,omitempty"`           // 14
	TLS           bool   `json:"tls,omitempty"`            // 15
	Transport     string `json:"transport,omitempty"`      // 16
	Path          string `json:"path,omitempty"`           // 17
	ServiceName   string `json:"service_name,omitempty"`   // 18
	Insecure      bool   `json:"insecure,omitempty"`       // 19
	SNI           string `json:"sni,omitempty"`            // 20
	Remark        string `json:"remark,omitempty"`         // 21
	ConnMode      string `json:"conn_mode,omitempty"`      // 22
	CountryCode   string `json:"country_code,omitempty"`   // 23
	Region        string `json:"region,omitempty"`         // 24
	Org           string `json:"org,omitempty"`            // 25
	VPN           string `json:"vpn,omitempty"`            // 26
}
