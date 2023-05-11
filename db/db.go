package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	conn *sql.DB
)

type DB struct {
	TotalAccount int
	uniqueIds    []string
	conn         *sql.DB
}

func New() *DB {
	db := DB{}
	if conn == nil {
		conn = db.connect()
	}

	db.conn = conn
	return &db
}

func (db *DB) connect() *sql.DB {
	connStr := os.Getenv("DB_URL")
	conn, _ := sql.Open("postgres", connStr)
	return conn
}

func (db *DB) isExists(values []any) bool {
	// Refer to data scheme
	// Server Port, UUID, Password, Transport, Conn Mode, VPN
	uid := fmt.Sprintf("%d_%s_%s_%s_%s_%s", values[2], values[3], values[4], values[16], values[22], values[26])

	if values[1] != "" {
		// Ip
		uid = uid + fmt.Sprintf("_%s", values[1])
	} else {
		if values[22] == "cdn" {
			// Host
			uid = uid + fmt.Sprintf("_%s", values[14])
		} else {
			// Server
			uid = uid + fmt.Sprintf("_%s", values[0])
		}
	}

	for _, existsId := range db.uniqueIds {
		if existsId == uid {
			return true
		}
	}

	db.uniqueIds = append(db.uniqueIds, uid)
	return false
}

func (db *DB) CreateTable() {
	query := `CREATE TABLE IF NOT EXISTS proxies (
		SERVER VARCHAR,
		IP VARCHAR,
		SERVER_PORT INTEGER,
		UUID VARCHAR,
		PASSWORD VARCHAR,
		SECURITY VARCHAR,
		ALTER_ID INTEGER,
		METHOD VARCHAR,
		PLUGIN VARCHAR,
		PLUGIN_OPTS VARCHAR,
		PROTOCOL VARCHAR,
		PROTOCOL_PARAM VARCHAR,
		OBFS VARCHAR,
		OBFS_PARAM VARCHAR,
		HOST VARCHAR,
		TLS INTEGER,
		TRANSPORT VARCHAR,
		PATH VARCHAR,
		SERVICE_NAME VARCHAR,
		INSECURE INTEGER,
		SNI VARCHAR,
		REMARK VARCHAR,
		CONN_MODE VARCHAR,
		COUNTRY_CODE VARCHAR,
		REGION VARCHAR,
		ORG VARCHAR,
		VPN VARCHAR
	)`

	if _, err := db.conn.Exec(query); err == nil {
		fmt.Println("[DB] Database successfully created !")
	} else {
		panic(err.Error())
	}
}
