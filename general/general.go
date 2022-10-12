package general

import (
	"github.com/trxbach/TOI-YEU-GO/db"
)

func InitTables(db *db.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS questions (id INTEGER PRIMARY KEY AUTO_INCREMENT, statement TEXT, answer TEXT)")
	if (err != nil) {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS accounts (id INTEGER PRIMARY KEY AUTO_INCREMENT, account TEXT, password TEXT)")
	if (err != nil) {
		return err
	}
	
	return nil
}