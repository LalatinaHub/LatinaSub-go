package db

// Merge from SS, SSR, Vmess, Vless/Trojan
type DBScheme struct {
	Server        string `json:"server,omitempty"`         // 0
	ServerPort    int    `json:"server_port"`              // 1
	UUID          string `json:"uuid,omitempty"`           // 2
	Password      string `json:"password,omitempty"`       // 3
	Security      string `json:"security"`                 // 4
	AlterId       int    `json:"alter_id,omitempty"`       // 5
	Method        string `json:"method,omitempty"`         // 6
	Plugin        string `json:"plugin,omitempty"`         // 7
	PluginOpts    string `json:"plugin_opts"`              // 8
	Protocol      string `json:"protocol,omitempty"`       // 9
	ProtocolParam string `json:"protocol_param,omitempty"` // 10
	OBFS          string `json:"obfs,omitempty"`           // 11
	OBFSParam     string `json:"obfs_param,omitempty"`     // 12
	Host          string `json:"host,omitempty"`           // 13
	TLS           bool   `json:"tls"`                      // 14
	Transport     string `json:"transport"`                // 15
	Path          string `json:"path,omitempty"`           // 16
	ServiceName   string `json:"service_name,omitempty"`   // 17
	Insecure      bool   `json:"insecure"`                 // 18
	SNI           string `json:"sni,omitempty"`            // 19
	Remark        string `json:"remark"`                   // 20
	ConnMode      string `json:"conn_mode"`                // 21
	CountryCode   string `json:"country_code"`             // 22
	Region        string `json:"region"`                   // 23
	VPN           string `json:"vpn"`                      // 24
}
