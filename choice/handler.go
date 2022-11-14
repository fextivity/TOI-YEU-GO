package choice

import (
	"github.com/trxbach/TOI-YEU-GO/database"
)

type Choice struct {
	Id int `json:"id"`
	Content string `json:"content"`
	Idq int `json:"idq"`
}

func InsertAnswerSql(db *database.DB, C *Choice) error {
	res, err := db.Exec("INSERT INTO choices (content, idq) VALUES (?, ?)", C.Content, C.Idq)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	C.Id = int(id)
	return nil
}