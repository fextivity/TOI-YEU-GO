package question

import (
	"github.com/trxbach/TOI-YEU-GO/choice"
	"github.com/trxbach/TOI-YEU-GO/database"
)

type Question struct {
	Id int `json:"id"`
	Content string `json:"content"`
	Idt int `json:"idt"`
	Choices []choice.Choice `json:"choices"`
}

func InsertQuestionSql(db *database.DB, Q *Question) error {
	res, err := db.Exec("INSERT INTO questions (content, idt) VALUES (?, ?)", Q.Content, Q.Idt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	Q.Id = int(id)
	for i := range Q.Choices {
		C := &Q.Choices[i]
		C.Idq = Q.Id
		err := choice.InsertAnswerSql(db, C)
		if err != nil {
			return err
		}
	}
	return nil
}

// func (wrp *Wrapper) AddQuestion(c echo.Context) error {
// 	var Q Question
// 	if err := c.Bind(&Q); err != nil {
// 		return c.String(http.StatusBadRequest, err.Error())
// 	}
// 	if err := InsertQuestionSql(wrp.db, &Q); err != nil {
// 		return c.String(http.StatusBadRequest, err.Error())
// 	}
// 	return c.JSON(http.StatusOK, Q)
// }