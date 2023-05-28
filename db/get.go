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
			server, ip, uuid, password, security, method, plugin, pluginOpts, protocol, protocolParam, obfs, obfsParam, host, transport, path, serviceName, sni, remark, connMode, countryCode, region, org, vpn sql.NullString
			id, serverPort, alterId                                                                                                                                                                              sql.NullInt32
			tls, insecure                                                                                                                                                                                        sql.NullBool
		)

		e := rows.Scan(
			&id,
			&server,
			&ip,
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

		if e != nil {
			fmt.Println(e)
		}

		result = append(result, DBScheme{
			Server:        server.String,
			Ip:            ip.String,
			ServerPort:    int(serverPort.Int32),
			UUID:          uuid.String,
			Password:      password.String,
			Security:      security.String,
			AlterId:       int(alterId.Int32),
			Method:        method.String,
			Plugin:        plugin.String,
			PluginOpts:    pluginOpts.String,
			Protocol:      protocol.String,
			ProtocolParam: protocolParam.String,
			OBFS:          obfs.String,
			OBFSParam:     obfsParam.String,
			Host:          host.String,
			TLS:           tls.Bool,
			Transport:     transport.String,
			Path:          path.String,
			ServiceName:   serviceName.String,
			Insecure:      insecure.Bool,
			SNI:           sni.String,
			Remark:        remark.String,
			ConnMode:      connMode.String,
			CountryCode:   countryCode.String,
			Region:        region.String,
			Org:           org.String,
			VPN:           vpn.String,
		})
	}

	return result
}
