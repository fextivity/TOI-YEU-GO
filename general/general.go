package general

import (
	"github.com/trxbach/TOI-YEU-GO/db"
)

func InitTables(db *db.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS accounts (id INTEGER PRIMARY KEY AUTO_INCREMENT, handle TEXT, password TEXT)")
	if (err != nil) {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS questions (id INTEGER PRIMARY KEY AUTO_INCREMENT, content TEXT)")
	if (err != nil) {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS answers (id INTEGER PRIMARY KEY AUTO_INCREMENT, content TEXT, idq INTEGER)")
	if (err != nil) {
		return err
	}
	
	return nil
}