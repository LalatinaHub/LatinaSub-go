package db

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/LalatinaHub/LatinaSub-go/account"
	"github.com/LalatinaHub/LatinaSub-go/sandbox"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func (db *DB) Save(boxes []*sandbox.SandBox) {
	var (
		values []string
	)

	db.TotalAccount = 0
	for _, box := range boxes {
		for _, value := range db.buildValuesQuery(box) {
			if value != "" {
				values = append(values, value)
				db.TotalAccount++
			}
		}
	}

	_, err := db.conn.Exec(fmt.Sprintf("INSERT INTO proxies VALUES %s", strings.Join(values, ", ")))
	if err != nil {
		fmt.Println("[DB] Failed to save accounts !")
		fmt.Println("[DB] Message:", err.Error())
		fmt.Println("[DB] Retrying ...")
		db.Save(boxes)
	} else {
		fmt.Println("[DB] Accounts saved !")
	}
}

func (db *DB) buildValuesQuery(box *sandbox.SandBox) []string {
	var (
		// Re-generate outbound to get pure config (without populated host)
		TLS         *option.OutboundTLSOptions
		Transport   *option.V2RayTransportOptions
		TLSSTR      string           = "NTLS"
		account     *account.Account = account.New(box.Link)
		anyOutbound interface{}
		values      []any
		queries     []string
	)

	switch account.Outbound.Type {
	case C.TypeVMess:
		anyOutbound = account.Outbound.VMessOptions
		TLS = account.Outbound.VMessOptions.TLS
		Transport = account.Outbound.VMessOptions.Transport
	case C.TypeTrojan:
		anyOutbound = account.Outbound.TrojanOptions
		TLS = account.Outbound.TrojanOptions.TLS
		Transport = account.Outbound.TrojanOptions.Transport
	case C.TypeVLESS:
		anyOutbound = account.Outbound.VLESSOptions
		TLS = account.Outbound.VLESSOptions.TLS
		Transport = account.Outbound.VLESSOptions.Transport
	case C.TypeShadowsocks:
		anyOutbound = account.Outbound.ShadowsocksOptions
	case C.TypeShadowsocksR:
		anyOutbound = account.Outbound.ShadowsocksROptions
	}

	// Null safe BEGIN
	if TLS == nil {
		TLS = &option.OutboundTLSOptions{}
	} else {
		if TLS.Enabled {
			TLSSTR = "TLS"
		}
	}
	if Transport == nil {
		Transport = &option.V2RayTransportOptions{
			Type:             "tcp", // unofficial, just for compatibility support
			WebsocketOptions: option.V2RayWebsocketOptions{},
			GRPCOptions:      option.V2RayGRPCOptions{},
			QUICOptions:      option.V2RayQUICOptions{},
			HTTPOptions:      option.V2RayHTTPOptions{},
		}
	} else {
		switch Transport.Type {
		case "ws", "grpc", "quic":
		case "websocket":
			Transport.Type = "ws"
		default:
			Transport.Type = "tcp"
		}
	}
	// Null safe END

	// Build values
	switch account.Outbound.Type {
	case C.TypeVMess:
		outbound := anyOutbound.(option.VMessOutboundOptions)
		values = []any{
			outbound.Server,
			outbound.ServerPort,
			outbound.UUID,
			"", // password
			outbound.Security,
			outbound.AlterId,
			"", // Method
			"", // Plugin
			"", // Plugin Opts
			"", // Protocol
			"", // Protocol Opts
			"", // OBFS
			"", // OBFS Param
		}
	case C.TypeTrojan:
		outbound := anyOutbound.(option.TrojanOutboundOptions)
		values = []any{
			outbound.Server,
			outbound.ServerPort,
			"", // UUID
			outbound.Password,
			"", // Security
			0,  // Alter ID
			"", // Method
			"", // Plugin
			"", // Plugin Opts
			"", // Protocol
			"", // Protocol Opts
			"", // OBFS
			"", // OBFS Param
		}
	case C.TypeVLESS:
		outbound := anyOutbound.(option.VLESSOutboundOptions)
		values = []any{
			outbound.Server,
			outbound.ServerPort,
			outbound.UUID,
			"", // Password
			"", // Security
			0,  // Alter ID
			"", // Method
			"", // Plugin
			"", // Plugin Opts
			"", // Protocol
			"", // Protocol Opts
			"", // OBFS
			"", // OBFS Param
		}
	case C.TypeShadowsocks:
		outbound := anyOutbound.(option.ShadowsocksOutboundOptions)
		values = []any{
			outbound.Server,
			outbound.ServerPort,
			"", // UUID
			outbound.Password,
			"", // Security
			0,  // Alter ID
			outbound.Method,
			outbound.Plugin,
			outbound.PluginOptions,
			"", // Protocol,
			"", // Protocol Opts
			"", // OBFS
			"", // OBFS Param
		}
	case C.TypeShadowsocksR:
		outbound := anyOutbound.(option.ShadowsocksROutboundOptions)
		values = []any{
			outbound.Server,
			outbound.ServerPort,
			"", // UUID
			outbound.Password,
			"", // Security
			0,  // Alter ID
			outbound.Method,
			"", // Plugin
			"", // Plugin Opts
			outbound.Protocol,
			outbound.ProtocolParam,
			outbound.Obfs,
			outbound.ObfsParam,
		}
	default:
		return queries
	}

	// Add TLS and Transport field to values
	values = append(values, []any{
		Transport.WebsocketOptions.Headers["Host"],
		TLS.Enabled,
		Transport.Type,
		Transport.WebsocketOptions.Path,
		Transport.GRPCOptions.ServiceName,
		TLS.Insecure,
		TLS.ServerName,
	}...)

	for _, mode := range box.ConnectMode {
		var valuesString string
		queryValues := append(values, []any{
			strings.ToUpper(fmt.Sprintf("%d %s %s %s %s %s", db.TotalAccount+1, box.IpapiObj.CountryCode, box.IpapiObj.Org, Transport.Type, mode, TLSSTR)),
			mode,
			box.IpapiObj.CountryCode,
			box.IpapiObj.Region,
			box.IpapiObj.Org,
			account.Outbound.Type,
		}...)

		// Check is account exists on database
		if db.isExists(queryValues) {
			continue
		}

		for _, value := range queryValues {
			switch reflect.TypeOf(value).Name() {
			case "bool":
				if value == true {
					valuesString = valuesString + `1, `
				} else {
					valuesString = valuesString + `0, `
				}
			case "int":
				valuesString = valuesString + fmt.Sprintf(`%d, `, value)
			case "uint16":
				valuesString = valuesString + fmt.Sprintf(`%d, `, value)
			default:
				valuesString = valuesString + fmt.Sprintf(`"%s", `, value)
			}
		}

		query := fmt.Sprintf(`(%s)`, strings.TrimSuffix(valuesString, ", "))
		queries = append(queries, query)
	}

	return queries
}
