package db

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/LalatinaHub/LatinaSub-go/account"
	"github.com/LalatinaHub/LatinaSub-go/helper"
	"github.com/LalatinaHub/LatinaSub-go/modules/bot"
	"github.com/LalatinaHub/LatinaSub-go/sandbox"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

var (
	ExcludedServer = func() []string {
		domains := []string{}
		rows, err := New().conn.Query("SELECT domain FROM domains;")
		if err != nil {
			fmt.Println(err)
			return domains
		}
		defer rows.Close()

		for rows.Next() {
			var (
				server sql.NullString
			)

			rows.Scan(&server)
			domains = append(domains, server.String)
		}

		return domains
	}()
)

func (db *DB) Save(boxes []*sandbox.SandBox) {
	var (
		values []string
		err    error
	)

	db.TotalAccount = 0
	for _, box := range boxes {
		for _, value := range db.BuildValuesQuery(box) {
			if value != "" {
				values = append(values, value)
				db.TotalAccount++
			}
		}
	}

	query := fmt.Sprintf(`INSERT INTO proxies (
		SERVER,
		IP,
		SERVER_PORT,
		UUID, PASSWORD,
		SECURITY,
		ALTER_ID,
		METHOD,
		PLUGIN,
		PLUGIN_OPTS,
		PROTOCOL,
		PROTOCOL_PARAM,
		OBFS,
		OBFS_PARAM,
		HOST,
		TLS,
		TRANSPORT,
		PATH,
		SERVICE_NAME,
		INSECURE,
		SNI,
		REMARK,
		CONN_MODE,
		COUNTRY_CODE,
		REGION,
		ORG,
		VPN
	) VALUES %s`, strings.ReplaceAll(strings.Join(values, ", "), `"`, "'"))

	for i := 0; i < 3; i++ {
		transactionQuery := fmt.Sprintf(`BEGIN; TRUNCATE proxies; %s; COMMIT;`, query)
		_, err = db.conn.Exec(transactionQuery)

		if err != nil {
			fmt.Println("[DB] Failed to save accounts !")
			fmt.Println("[DB] Message:", err.Error())
			fmt.Println("[DB] Retrying ...")
		} else {
			fmt.Println("[DB] Accounts saved !")
			break
		}

		if i >= 2 {
			fmt.Println("[DB] Retry attempt exceeded !")
			f, _ := os.Create("DB_QUERY.txt")
			defer f.Close()

			f.WriteString(transactionQuery)

			fmt.Println("Sending database query to admin...")
			bot.SendTextFileToAdmin("db_query.txt", transactionQuery)
		}
	}
}

func (db *DB) BuildValuesQuery(box *sandbox.SandBox) []string {
	var (
		// Re-generate outbound to get pure config (without populated host)
		TLS         *option.OutboundTLSOptions
		Transport   *option.V2RayTransportOptions
		TLSSTR      string           = "NTLS"
		host        string           = ""
		account     *account.Account = account.New(box.Outbound)
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

		if m, _ := regexp.MatchString("tls", account.Outbound.ShadowsocksROptions.Obfs); m {
			TLS = &option.OutboundTLSOptions{
				Enabled: true,
			}
		}
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
		if len(Transport.WebsocketOptions.Headers["Host"]) > 0 {
			host = Transport.WebsocketOptions.Headers["Host"][0]
		}

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
			box.Geoip.Ip,
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
			box.Geoip.Ip,
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
			box.Geoip.Ip,
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
			box.Geoip.Ip,
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
			box.Geoip.Ip,
			outbound.ServerPort,
			"", // UUID
			outbound.Password,
			"", // Security
			0,  // Alter ID
			outbound.Method,
			"", // Plugin
			"", // Plugin Options
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
		host,
		TLS.Enabled,
		Transport.Type,
		Transport.WebsocketOptions.Path,
		Transport.GRPCOptions.ServiceName,
		TLS.Insecure,
		host,
	}...)

	for _, mode := range box.ConnectMode {
		var valuesString string
		queryValues := append(values, []any{
			strings.ToUpper(fmt.Sprintf("%d %s %s %s %s %s", db.TotalAccount+len(queries)+1, helper.CCToEmoji(box.Geoip.CountryCode), box.Geoip.Org, Transport.Type, mode, TLSSTR)),
			mode,
			box.Geoip.CountryCode,
			box.Geoip.Region,
			box.Geoip.Org,
			account.Outbound.Type,
		}...)

		// Check if account server is excluded
		for _, server := range ExcludedServer {
			if server == queryValues[0] || server == queryValues[14] {
				return queries
			}
		}

		// Check if account exists on database
		if db.isExists(queryValues) {
			continue
		}

		for _, value := range queryValues {
			switch reflect.TypeOf(value).Kind() {
			case reflect.Bool:
				if value == true {
					valuesString = valuesString + `1, `
				} else {
					valuesString = valuesString + `0, `
				}
			case reflect.String:
				valuesString = valuesString + fmt.Sprintf(`"%s", `, value)
			default:
				valuesString = valuesString + fmt.Sprintf(`%d, `, value)
			}
		}

		query := fmt.Sprintf(`(%s)`, strings.TrimSuffix(valuesString, ", "))
		queries = append(queries, query)
	}

	if len(box.ConnectMode) <= 0 {
		var tempValues []string

		for _, value := range values {
			tempValues = append(tempValues, fmt.Sprintf("%v", value))
		}
		tempValues = append(tempValues, account.Outbound.Type)

		return tempValues
	}

	return queries
}
