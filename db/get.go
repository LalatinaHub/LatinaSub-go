package db

import (
	"database/sql"
	"fmt"
)

func (db *DB) Get(filter string) []DBScheme {
	query := fmt.Sprintf(`SELECT * FROM proxies %s`, filter)
	rows, err := db.conn.Query(query)
	if err != nil {
		fmt.Println(err)

		return []DBScheme{}
	}
	defer rows.Close()

	return toJson(rows)
}

func toJson(rows *sql.Rows) []DBScheme {
	var result []DBScheme

	if rows == nil {
		return result
	}

	for rows.Next() {
		var (
			server, uuid, password, security, method, plugin, pluginOpts, protocol, protocolParam, obfs, obfsParam, host, transport, path, serviceName, sni, remark, connMode, countryCode, region, org, vpn string
			serverPort, alterId                                                                                                                                                                              int
			tls, insecure                                                                                                                                                                                    bool
		)

		rows.Scan(
			&server,
			&serverPort,
			&uuid,
			&password,
			&security,
			&alterId,
			&method,
			&plugin,
			&pluginOpts,
			&protocol,
			&protocolParam,
			&obfs,
			&obfsParam,
			&host,
			&tls,
			&transport,
			&path,
			&serviceName,
			&insecure,
			&sni,
			&remark,
			&connMode,
			&countryCode,
			&region,
			&org,
			&vpn)

		result = append(result, DBScheme{
			Server:        server,
			ServerPort:    serverPort,
			UUID:          uuid,
			Password:      password,
			Security:      security,
			AlterId:       alterId,
			Method:        method,
			Plugin:        plugin,
			PluginOpts:    pluginOpts,
			Protocol:      protocol,
			ProtocolParam: protocolParam,
			OBFS:          obfs,
			OBFSParam:     obfsParam,
			Host:          host,
			TLS:           tls,
			Transport:     transport,
			Path:          path,
			ServiceName:   serviceName,
			Insecure:      insecure,
			SNI:           sni,
			Remark:        remark,
			ConnMode:      connMode,
			CountryCode:   countryCode,
			Region:        region,
			Org:           org,
			VPN:           vpn,
		})
	}

	return result
}
