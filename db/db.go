package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
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
		ORG VARCHAR NOT NULL ON CONFLICT REPLACE DEFAULT "VPS",
		VPN VARCHAR
	)`

	if _, err := db.conn.Exec(query); err == nil {
		fmt.Println("[DB] Database successfully created !")
	} else {
		panic(err.Error())
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
