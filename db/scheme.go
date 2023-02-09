package db

// Merge from SS, SSR, Vmess, Vless/Trojan
type DBScheme struct {
	Server        string `json:"server,omitempty"`         // 0
	ServerPort    int    `json:"server_port,omitempty"`    // 1
	UUID          string `json:"uuid,omitempty"`           // 2
	Password      string `json:"password,omitempty"`       // 3
	Security      string `json:"security,omitempty"`       // 4
	AlterId       int    `json:"alter_id,omitempty"`       // 5
	Method        string `json:"method,omitempty"`         // 6
	Plugin        string `json:"plugin,omitempty"`         // 7
	PluginOpts    string `json:"plugin_opts,omitempty"`    // 8
	Protocol      string `json:"protocol,omitempty"`       // 9
	ProtocolParam string `json:"protocol_param,omitempty"` // 10
	OBFS          string `json:"obfs,omitempty"`           // 11
	OBFSParam     string `json:"obfs_param,omitempty"`     // 12
	Host          string `json:"host,omitempty"`           // 13
	TLS           bool   `json:"tls,omitempty"`            // 14
	Transport     string `json:"transport,omitempty"`      // 15
	Path          string `json:"path,omitempty"`           // 16
	ServiceName   string `json:"service_name,omitempty"`   // 17
	Insecure      bool   `json:"insecure,omitempty"`       // 18
	SNI           string `json:"sni,omitempty"`            // 19
	Remark        string `json:"remark,omitempty"`         // 20
	ConnMode      string `json:"conn_mode,omitempty"`      // 21
	CountryCode   string `json:"country_code,omitempty"`   // 22
	Region        string `json:"region,omitempty"`         // 23
	Org           string `json:"org,omitempty"`            // 24
	VPN           string `json:"vpn,omitempty"`            // 25
}
