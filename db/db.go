package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/LalatinaHub/LatinaSub-go/account"
	"github.com/LalatinaHub/LatinaSub-go/sandbox"
	_ "github.com/mattn/go-sqlite3"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

var (
	DbPath string = "result/"
	DbName string = "db.sqlite"
	DbFile string = DbPath + DbName
)

type DB struct {
	TotalAccount int
	uniqueIds    []string
	conn         *sql.DB
}

func New() *DB {
	db := DB{}
	db.conn = db.connect()

	return &db
}

func (db *DB) connect() *sql.DB {
	conn, _ := sql.Open("sqlite3", DbFile)
	return conn
}

func (db *DB) isExists(values []any) bool {
	// Refer to data scheme
	// Server Port, UUID, Password, Transport, Path, Service Name, Conn Mode, VPN
	id := fmt.Sprintf(`%d_"%s"_"%s"_"%s"_"%s"_"%s"_"%s"_"%s"`, values[1], values[2], values[3], values[15], values[16], values[17], values[21], values[24])

	if values[23] == "cdn" {
		// Host
		id = id + fmt.Sprintf(`_"%s"`, values[13])
	} else {
		// Server
		id = id + fmt.Sprintf(`_"%s"`, values[0])
	}

	for _, existsId := range db.uniqueIds {
		if existsId == id {
			return true
		}
	}

	db.uniqueIds = append(db.uniqueIds, id)
	return false
}

func (db *DB) FlushAndCreate() {
	// Remove previous database
	if file, _ := os.Stat(DbFile); file != nil {
		os.Remove(DbFile)
	}

	query := `CREATE TABLE IF NOT EXISTS proxies (
		SERVER VARCHAR,
		SERVER_PORT INTEGER,
		UUID VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "",
		PASSWORD VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "",
		SECURITY VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "",
		ALTER_ID INTEGER NOT NULL ON CONFLICT REPLACE DEFAULT 0,
		METHOD VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "",
		PLUGIN VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "",
		PLUGIN_OPTS VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "",
		PROTOCOL VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "",
		PROTOCOL_PARAM VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "",
		OBFS VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "",
		OBFS_PARAM VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "",
		HOST VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "",
		TLS INTEGER NOT NULL ON CONFLICT REPLACE DEFAULT 1,
		TRANSPORT VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "",
		PATH VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "",
		SERVICE_NAME VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "",
		INSECURE INTEGER NOT NULL ON CONFLICT REPLACE DEFAULT 1,
		SNI VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "",
		REMARK VARCHAR,
		CONN_MODE VARCHAR,
		COUNTRY_CODE VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "XX",
		REGION VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "Unknown",
		VPN VARCHAR
	)`

	if _, err := db.conn.Exec(query); err == nil {
		fmt.Println("[DB] Database successfully created !")
	} else {
		panic(err.Error())
	}
}

func (db *DB) Save(box *sandbox.SandBox) {
	var (
		// Re-generate outbound to get pure config (without populated host)
		TLS         *option.OutboundTLSOptions
		Transport   *option.V2RayTransportOptions
		account     *account.Account = account.New(box.Link)
		anyOutbound interface{}
		values      []any
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
	}
	if Transport == nil {
		Transport = &option.V2RayTransportOptions{
			Type:             "tcp", // unofficial, just for compatibility support
			WebsocketOptions: option.V2RayWebsocketOptions{},
			GRPCOptions:      option.V2RayGRPCOptions{},
			QUICOptions:      option.V2RayQUICOptions{},
			HTTPOptions:      option.V2RayHTTPOptions{},
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
			"", // Alter ID
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
			"", // Alter ID
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
			"", // Alter ID
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
			"", // Alter ID
			outbound.Method,
			"", // Plugin
			"", // Plugin Opts
			outbound.Protocol,
			outbound.ProtocolParam,
			outbound.Obfs,
			outbound.ObfsParam,
		}
	default:
		return
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
			account.Outbound.Tag,
			mode,
			box.IpapiObj.CountryCode,
			box.IpapiObj.Region,
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

		query := fmt.Sprintf(`INSERT INTO proxies VALUES (%s)`, strings.TrimSuffix(valuesString, ", "))

		_, err := db.conn.Exec(query)
		if err != nil {
			// Force trying to insert account
			if err.Error() == "database is locked" {
				fmt.Println("[DB] Database locked, retrying ...")
				time.Sleep(100 * time.Millisecond)
				db.Save(box)
			} else {
				panic(err)
			}
		} else {
			db.TotalAccount++
		}
	}
}

func Init() {
	// Check and create dir "result/"
	if _, err := os.Stat(DbPath); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(DbPath, os.ModePerm)
		} else {
			log.Panic(err)
		}
	}
}
