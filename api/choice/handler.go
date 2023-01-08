package choice

import (
	"github.com/trxbach/TOI-YEU-GO/database"
)

type Choice struct {
	Id int `json:"id"`
	Content string `json:"content"`
	IsAnswer bool `json:"is_answer"`
	Idq int `json:"idq"`
}

func InsertAnswerSql(db *database.DB, C *Choice) error {
	res, err := db.Exec(`INSERT INTO choices (content, is_answer, idq) VALUES (?, ?, ?)`, C.Content, C.IsAnswer, C.Idq)
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

func DeleteChoiceSQL(db *database.DB, id int) error {
	_, err := db.Exec(`DELETE FROM questions WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return nil
}