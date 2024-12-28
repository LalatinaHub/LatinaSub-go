package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var conn *sql.DB

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

func (db *DB) Conn() *sql.DB {
	return db.conn
}

func (db *DB) connect() *sql.DB {
	connStr := os.Getenv("DB_URL")
	conn, _ := sql.Open("postgres", connStr)
	return conn
}

func (db *DB) isExists(values []any) bool {
	// Refer to data scheme
	// Server Port, UUID, Password, OBFS Param, Path, Transport, Conn Mode, VPN
	uid := fmt.Sprintf("%d_%s_%s_%s_%s_%s_%s_%s", values[2], values[3], values[4], values[13], values[17], values[16], values[22], values[26])

	if values[22] == "cdn" {
		// Host
		uid = uid + fmt.Sprintf("_%s", values[14])
	} else {
		// Server
		uid = uid + fmt.Sprintf("_%s", values[0])
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
		ID SERIAL PRIMARY KEY,
		SERVER TEXT,
		IP TEXT,
		SERVER_PORT INT8,
		UUID TEXT,
		PASSWORD TEXT,
		SECURITY TEXT,
		ALTER_ID INT2,
		METHOD TEXT,
		PLUGIN TEXT,
		PLUGIN_OPTS TEXT,
		PROTOCOL TEXT,
		PROTOCOL_PARAM TEXT,
		OBFS TEXT,
		OBFS_PARAM TEXT,
		HOST TEXT,
		TLS INT2,
		TRANSPORT TEXT,
		PATH TEXT,
		SERVICE_NAME TEXT,
		INSECURE INT2,
		SNI TEXT,
		REMARK TEXT,
		CONN_MODE TEXT,
		COUNTRY_CODE TEXT,
		REGION TEXT,
		ORG TEXT,
		VPN TEXT
	)`

	if _, err := db.conn.Exec(query); err == nil {
		fmt.Println("[DB] Database successfully created !")
	} else {
		panic(err.Error())
	}
}
