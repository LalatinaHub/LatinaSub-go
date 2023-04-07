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
	id := fmt.Sprintf(`%d_"%s"_"%s"_"%s"_"%s"_"%s"`, values[1], values[2], values[3], values[15], values[21], values[25])

	if values[21] == "cdn" {
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

func (db *DB) CreateTable() {
	query := `CREATE TABLE IF NOT EXISTS proxies (
		SERVER VARCHAR,
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
