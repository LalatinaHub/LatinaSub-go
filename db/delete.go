package db

import "fmt"

func (db *DB) Delete(remark string) {
	query := fmt.Sprintf(`DELETE FROM proxies WHERE REMARK='%s'`, remark)
	_, _ = db.conn.Exec(query)
}
